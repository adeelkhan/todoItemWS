import axios from "axios";
import React, { useState, useContext } from "react";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "../AuthContext";

const Login = () => {
  const [userName, setUserName] = useState("");
  const [userPassword, setUserPassword] = useState("");
  const navigate = useNavigate();

  const { setAuthUser } = useContext(AuthContext);

  const editUser = (username) => {
    setUserName(username);
  };
  const editPassword = (password) => {
    setUserPassword(password);
  };

  const LoginUser = () => {
    axios
      .post(
        "http://localhost:8090/signin",
        {
          username: userName,
          password: userPassword,
        },
        {
          withCredentials: true,
        }
      )
      .then((response) => {
        setAuthUser(response.data.user);
        navigate("/list");
      });
  };

  return (
    <>
      User :
      <input
        type="text"
        value={userName}
        label="username"
        onChange={(e) => editUser(e.target.value)}
      />
      Password :
      <input
        type="password"
        value={userPassword}
        label="password"
        onChange={(e) => editPassword(e.target.value)}
      />
      <button onClick={LoginUser}>Login</button>
    </>
  );
};
export default Login;
