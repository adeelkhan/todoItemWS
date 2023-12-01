import logo from "./logo.svg";
import "./App.css";
// import CreateItem from "./components/CreateItem";
// import ListItem from "./components/ListItem";
// import DeleteItem from "./components/DeleteItem";
// import UpdateItem from "./components/UpdateItem";
import { useState, useEffect } from "react";
import axios from "axios";
import { BrowserRouter, Routes, Route } from "react-router-dom";

function App() {
  const [todoItem, setTodoItem] = useState("");
  const [items, setItems] = useState([]);
  const [itemUpdated, setItemUpdated] = useState(false);
  const [updateBox, setShowUpdateBox] = useState([]);
  const [newItemValue, setNewItemValue] = useState("");

  const editTodoItem = (e) => {
    setTodoItem(e.target.value);
  };

  const createTodoItem = () => {
    if (todoItem == "") {
      return;
    }
    axios
      .post("http://localhost:8090/create", {
        item_name: todoItem,
      })
      .then((response) => {
        setItemUpdated(true);
        setTodoItem("");
      });
  };

  useEffect(() => {
    axios.get("http://localhost:8090/list").then((data) => {
      data = data.data;
      data.items.sort((a, b) => {
        return a.Id - b.Id;
      });
      setItems(data.items);
      setItemUpdated(false);
    });
  }, [itemUpdated]);

  const deleteTodoItem = (Id) => {
    axios
      .post("http://localhost:8090/delete", {
        item_id: Id,
      })
      .then((response) => {
        setItemUpdated(true);
      });
  };
  const updateTodoItem = (Id, name) => {
    axios
      .post("http://localhost:8090/update", {
        item_id: Id,
        item_name: newItemValue,
      })
      .then((response) => {
        setNewItemValue("");
        removeUpdateBox(Id);
        setItemUpdated(true);
      });
  };

  const setUpdateBox = (Id) => {
    let updateList = [...updateBox];
    let found = updateList.includes(Id);
    if (!found) {
      updateList.push(Id);
      setShowUpdateBox(updateList);
    }
  };
  const removeUpdateBox = (Id) => {
    let updateList = [...updateBox].filter((e) => e != Id);
    setShowUpdateBox(updateList);
  };

  const IsItemEditable = (Id) => {
    return updateBox.includes(Id);
  };

  const editNewItemValue = (e) => {
    setNewItemValue(e.target.value);
  };

  return (
    <div className="App">
      <div>Todoapp</div>
      <div>
        <input value={todoItem} onChange={(e) => editTodoItem(e)} />
        <button onClick={() => createTodoItem()}>Add Item</button>
      </div>
      <ul>
        {items &&
          items.map((item) => {
            return (
              <li key={item.Id} onClick={() => setUpdateBox(item.Id)}>
                {item.item_name}
                <button onClick={() => deleteTodoItem(item.Id)}>X</button>
                {IsItemEditable(item.Id) && (
                  <>
                    <input
                      value={newItemValue}
                      onChange={(e) => editNewItemValue(e)}
                    />
                    <button
                      onClick={() => updateTodoItem(item.Id, item.item_name)}
                    >
                      Edit
                    </button>
                    <button onClick={() => removeUpdateBox(item.Id)}>
                      Close
                    </button>
                  </>
                )}
              </li>
            );
          })}
      </ul>
    </div>
  );
}

export default App;
