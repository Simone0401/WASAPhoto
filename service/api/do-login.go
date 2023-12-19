package api

import (
	"encoding/json"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/Simone0401/WASAPhoto/service/database"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// doLogin is the handler for the API endpoint POST /session.
// It takes the username from the request body and returns the user object and the authorization token in a JSON object.
// If the user does not exist, it creates a new user.
// The request body must be a JSON object with the following fields:
//   - username: string
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	var user User

	// trying parsing request object to API User Struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		context.Logger.Error(err.Error())
		http.Error(w, "Error parsing JSON Object request body", http.StatusBadRequest)
		return
	}

	// check if the username is valid
	if !user.IsValid() {
		http.Error(w, "Username format is not valid!", http.StatusBadRequest)
		return
	}

	// check if the user already exists
	// if not exists we need creating that
	// otherwise return user object
	var userdb database.User
	userdb = user.ToDatabase()
	exists, err := rt.db.CheckExistsByUsername(userdb.Username)
	if err != nil {
		context.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	if !exists {
		// create the user
		user, err = createUser(rt, user)
		if err != nil {
			context.Logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		context.Logger.Info("User correctly created. ", user)
		w.WriteHeader(http.StatusCreated)

	} else {
		// recover the user
		user, err = recoverUserByUsername(rt, user)
		if err != nil {
			context.Logger.Info(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		context.Logger.Info("User correctly recovered!", user)
		w.WriteHeader(http.StatusOK)

	}

	_ = json.NewEncoder(w).Encode(user)

}

func createUser(rt *_router, user User) (User, error) {
	var userdb database.User

	userdb = user.ToDatabase()
	userdb, err := rt.db.CreateUser(userdb.Username)

	if err != nil {
		return User{}, err
	}

	err = user.FromDatabase(userdb)

	if err != nil {
		return User{}, err
	}

	return user, err
}

func recoverUserByUsername(rt *_router, user User) (User, error) {
	var userdb database.User

	userdb = user.ToDatabase()
	userdb, err := rt.db.GetUserByUsername(user.Username)

	if err != nil {
		return User{}, err
	}

	err = user.FromDatabase(userdb)

	if err != nil {
		return User{}, err
	}

	return user, err
}
