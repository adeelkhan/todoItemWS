import { useState } from "react";
import axios from "axios";

const CreateItem = () => {
  const [todoItem, setTodoItem] = useState();

  const editTodoItem = (e) => {
    setTodoItem(e.target.value);
  };

  const createTodoItem = () => {
    axios.post("http://localhost:8090/create", {
      item_name: todoItem,
    });
  };

  return (
    <div>
      <input value={todoItem} onChange={(e) => editTodoItem(e)} />
      <button onClick={createTodoItem}>Add Item</button>
    </div>
  );
};

export default CreateItem;
