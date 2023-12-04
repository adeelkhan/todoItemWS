import React, { useContext, useEffect } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "../AuthContext";

const LogOut = () => {
  const navigate = useNavigate();
  const { setAuthUser } = useContext(AuthContext);
  useEffect(() => {
    axios.get("http://localhost:8090/logout").then(() => {
      setAuthUser("");
      navigate("/login");
    });
  }, []);
};

export default LogOut;
