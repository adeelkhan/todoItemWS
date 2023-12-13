import logo from "../logo.svg";
import "../App.css";
import { useState, useEffect, useContext } from "react";
import axios from "axios";
import { AuthContext } from "../AuthContext";
import { useNavigate } from "react-router-dom";
import showErrorBanner from "./Error";

function App() {
  const [todoItem, setTodoItem] = useState("");
  const [items, setItems] = useState([]);
  const [itemUpdated, setItemUpdated] = useState(false);
  const [updateBox, setShowUpdateBox] = useState("");
  const [newItemValue, setNewItemValue] = useState("");
  const { username } = useContext(AuthContext);
  const [error, setError] = useState();
  const navigate = useNavigate();

  const editTodoItem = (e) => {
    setTodoItem(e.target.value);
  };

  const processError = (response) => {
    if (response.status == 401) {
      navigate("/login");
    } else if (response.status === 400) {
      setError("Something unexpected happened");
    }
  };
  const createTodoItem = () => {
    if (todoItem === "") {
      return;
    }
    axios
      .post(
        "http://localhost:8090/create",
        {
          item_name: todoItem,
        },
        {
          withCredentials: true,
        }
      )
      .then((response) => {
        setItemUpdated(true);
        setTodoItem("");
      })
      .catch(({ response }) => {
        processError(response);
      });
  };

  useEffect(() => {
    axios
      .get("http://localhost:8090/list", {
        withCredentials: true,
      })
      .then((response) => {
        let data = response.data;
        setItems(data.items);
        setItemUpdated(false);
      })
      .catch(({ response }) => {
        processError(response);
      });
  }, [itemUpdated]);

  const deleteTodoItem = (Id) => {
    axios
      .post(
        "http://localhost:8090/delete",
        {
          item_id: Id,
        },
        {
          withCredentials: true,
        }
      )
      .then((response) => {
        setItemUpdated(true);
      })
      .catch(({ response }) => {
        processError(response);
      });
  };
  const updateTodoItem = (Id, name) => {
    if (newItemValue == "") {
      alert("Todo item cannot be empty");
      return;
    }

    axios
      .post(
        "http://localhost:8090/update",
        {
          item_id: Id,
          item_name: newItemValue,
        },
        {
          withCredentials: true,
        }
      )
      .then((response) => {
        setNewItemValue("");
        setShowUpdateBox("");
        setItemUpdated(true);
      })
      .catch(({ response }) => {
        processError(response);
      });
  };

  const IsItemEditable = (Id) => {
    return updateBox && updateBox === Id;
  };

  const editNewItemValue = (e) => {
    setNewItemValue(e.target.value);
  };
  const logOut = (e) => {
    navigate("/logout");
  };

  return (
    <div>
      <header>
        <div className="flex justify-between bg-slate-900 h-10 items-center">
          <div>
            <img
              className="mx-auto h-10 w-auto"
              src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=600"
              alt="Your Company"
            />
          </div>
          <div>Todo App</div>
          <div>
            <a
              href=""
              onClick={logOut}
              className="lg:inline-block py-2 px-6 bg-blue-500 hover:bg-blue-600 text-sm text-white font-bold rounded-xl transition duration-200"
            >
              Logout
            </a>
          </div>
        </div>
      </header>
      <main>
        {error && showErrorBanner(error)}
        <div>
          <input
            value={todoItem}
            onChange={(e) => editTodoItem(e)}
            className="rounded-xl py-1"
          />
          <button
            onClick={() => createTodoItem()}
            className="lg:inline-block py-1 px-2 bg-blue-500 hover:bg-blue-600 text-sm text-white font-bold rounded-xl transition duration-200"
          >
            Add Item
          </button>
        </div>
        <ul>
          {items &&
            items.map((item) => {
              return (
                <li key={item.Id}>
                  <button
                    onClick={() => deleteTodoItem(item.Id)}
                    className="lg:inline-block py-1 px-1 bg-blue-500 hover:bg-blue-600 text-sm text-white font-bold rounded-xl transition duration-200"
                  >
                    X
                  </button>
                  <span onClick={() => setShowUpdateBox(item.Id)}>
                    {item.item_name}
                  </span>

                  {IsItemEditable(item.Id) && (
                    <>
                      <input
                        value={newItemValue}
                        onChange={(e) => editNewItemValue(e)}
                        className="rounded-xl py-1"
                      />
                      <button
                        onClick={() => updateTodoItem(item.Id, item.item_name)}
                        className="lg:inline-block py-1 px-2 bg-blue-500 hover:bg-blue-600 text-sm text-white font-bold rounded-xl transition duration-200"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => setShowUpdateBox("")}
                        className="lg:inline-block py-1 px-2 bg-blue-500 hover:bg-blue-600 text-sm text-white font-bold rounded-xl transition duration-200"
                      >
                        Close
                      </button>
                    </>
                  )}
                </li>
              );
            })}
        </ul>
      </main>
      <footer></footer>
    </div>
  );
}

export default App;
