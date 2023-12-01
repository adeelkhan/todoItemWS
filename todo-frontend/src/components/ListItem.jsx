import { useEffect, useState } from "react";
import CreateItem from "./CreateItem";
import axios from "axios";

const ListItem = () => {
  const [items, setItems] = useState([]);
  const [itemUpdated, setItemUpdated] = useState(Boolean);
  useEffect(() => {
    axios.get("http://localhost:8090/list").then((data) => {
      data = data.data;
      setItems(data.items);
    });
    setItemUpdated(false);
  }, []);

  return (
    <div>
      <CreateItem />
      <p>Todo Items:</p>
      <ul>
        {items &&
          items.map((item) => {
            return <li key={item.Id}>{item.item_name}</li>;
          })}
      </ul>
    </div>
  );
};

export default ListItem;
