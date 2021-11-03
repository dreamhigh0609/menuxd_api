package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/super-radmir/menuxd_api/pkg/dish"
	"github.com/super-radmir/menuxd_api/pkg/notification"
	"github.com/super-radmir/menuxd_api/pkg/order"
	"github.com/super-radmir/menuxd_api/pkg/table"
	melody "gopkg.in/olahol/melody.v1"
)

var tableTypes map[string]string = map[string]string{"table": "Mesa", "bar": "Bar"}

// OrderRouter is a router to orders.
type OrderRouter struct {
	OrderStorage  order.Storage
	TableStorage  table.Storage
	DishStorage   dish.Storage
	MessageStream chan notification.Notification
}

func ordersHandler(m *melody.Melody) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := m.HandleRequest(w, r)
		if err != nil {
			http.Error(w, "Websocket connection failed", http.StatusInternalServerError)
			return
		}
	}
}

func (or OrderRouter) callWaiter(w http.ResponseWriter, r *http.Request) {
	tableIDStr := chi.URLParam(r, "tableId")
	tableID, err := strconv.Atoi(tableIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t, err := or.TableStorage.GetByID(uint(tableID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	n := notification.Notification{
		Type:     notification.CallWaiter,
		Message:  fmt.Sprintf("%v, #%d solicita al mozo", tableTypes[t.Type], t.Number),
		Date:     time.Now(),
		ClientID: uint(t.ClientID),
		Active:   true,
		Table:    &t,
	}

	go func() {
		or.MessageStream <- n
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (or OrderRouter) getBill(w http.ResponseWriter, r *http.Request) {
	tableIDStr := chi.URLParam(r, "tableId")
	tableID, err := strconv.Atoi(tableIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t, err := or.TableStorage.GetByID(uint(tableID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	n := notification.Notification{
		Type:     notification.GetCheck,
		Message:  fmt.Sprintf("%v #%d solicita la cuenta", tableTypes[t.Type], t.Number),
		Date:     time.Now(),
		ClientID: t.ClientID,
		Active:   true,
		Table:    &t,
	}

	go func() {
		or.MessageStream <- n
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (or OrderRouter) createHandler(w http.ResponseWriter, r *http.Request) {
	o := order.Order{}
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	o, err = or.OrderStorage.Create(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	j, err := json.Marshal(o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func (or OrderRouter) addItemHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clientIDStr := chi.URLParam(r, "clientId")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	items := []order.Item{}
	err = json.NewDecoder(r.Body).Decode(&items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	err = or.OrderStorage.Add(uint(id), items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	n := notification.Notification{}

	go func() {

		for _, i := range items {
			var pic string
			if i.Dish != nil {
				i.DishID = i.Dish.ID
			}
			storedDish, err := or.DishStorage.GetByID(i.DishID)
			if err == nil {
				pic = storedDish.Pictures[0]
			}
			n.Type = notification.MakeOrder
			n.Picture = pic
			n.ClientID = uint(clientID)
			n.Active = true
			n.Date = time.Now()
			msg := fmt.Sprintf("Orden recibida, %s", storedDish.Name)

			storedOrder, err := or.OrderStorage.GetByID(i.OrderID)
			if err == nil {
				msg += fmt.Sprintf(", %s #%d", storedOrder.Table.Type, storedOrder.Table.Number)
			}
			n.Table = storedOrder.Table

			n.Message = msg

			or.MessageStream <- n
		}
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (or OrderRouter) updateHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	o := order.Order{}
	err = json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	err = or.OrderStorage.Update(uint(id), &o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (or OrderRouter) updateItemHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	err = or.OrderStorage.PatchItem(uint(id), m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (or OrderRouter) getAllHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientId")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orders, err := or.OrderStorage.GetAll(uint(clientID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	j, err := json.Marshal(orders)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (or OrderRouter) getActiveHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientId")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orders, err := or.OrderStorage.GetAllActive(uint(clientID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	j, err := json.Marshal(orders)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// NewOrderRouter returns the order's handler with default configuration.
func NewOrderRouter(s order.Storage, ts table.Storage, ds dish.Storage) *chi.Mux {
	ch := make(chan notification.Notification, 100)
	or := OrderRouter{
		OrderStorage:  s,
		TableStorage:  ts,
		DishStorage:   ds,
		MessageStream: ch,
	}

	r := chi.NewRouter()

	m := melody.New()
	r.Get("/{clientId}/ws", ordersHandler(m))

	m.HandleMessage(messageHandler(m, ch))
	m.HandleConnect(connectHandler(m))

	go func() {

		for {
			n := <-or.MessageStream
			j, err := json.Marshal(n)
			if err != nil {
				m.CloseWithMsg([]byte(err.Error()))
				return
			}

			m.BroadcastFilter(j, func(ss *melody.Session) bool {
				clientIDStr := chi.URLParam(ss.Request, "clientId")
				clientID, err := strconv.Atoi(clientIDStr)
				if err != nil {
					return false
				}

				return uint(clientID) == n.ClientID
			})
		}
	}()

	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("XD_SIGNING_STRING")), nil)

	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Get("/call/{tableId}/waiter", or.callWaiter)
	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Get("/call/{tableId}/bill", or.getBill)
	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Post("/", or.createHandler)
	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Put("/{id}", or.updateHandler)
	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Put("/add/{id}/client/{clientId}", or.addItemHandler)
	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Patch("/item/{id}", or.updateItemHandler)
	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Get("/client/{clientId}", or.getAllHandler)
	r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator).
		Get("/active/{clientId}", or.getActiveHandler)

	return r
}

func messageHandler(m *melody.Melody, MessageStream chan notification.Notification) func(*melody.Session, []byte) {
	return func(s *melody.Session, _ []byte) {

		n := <-MessageStream
		j, err := json.Marshal(n)
		if err != nil {
			s.CloseWithMsg([]byte(err.Error()))
			return
		}

		m.BroadcastFilter(j, func(ss *melody.Session) bool {
			clientIDStr := chi.URLParam(ss.Request, "clientId")
			clientID, err := strconv.Atoi(clientIDStr)
			if err != nil {
				return false
			}

			return uint(clientID) == n.ClientID
		})

	}
}

func connectHandler(m *melody.Melody) func(*melody.Session) {
	return func(s *melody.Session) {
		clientIDStr := chi.URLParam(s.Request, "clientId")
		clientID, err := strconv.Atoi(clientIDStr)
		if err != nil {
			s.CloseWithMsg([]byte(err.Error()))
			return
		}

		n := notification.Notification{
			Type:     notification.Connected,
			Message:  "Connected",
			Date:     time.Now(),
			ClientID: uint(clientID),
			Active:   true,
		}

		j, err := json.Marshal(n)
		if err != nil {
			s.CloseWithMsg([]byte(err.Error()))
			return
		}

		m.BroadcastFilter(j, func(ss *melody.Session) bool {
			msgClientID := chi.URLParam(ss.Request, "clientId")
			return clientIDStr == msgClientID
		})
	}
}
