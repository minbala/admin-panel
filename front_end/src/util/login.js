import { message } from "antd";
import { endpoints } from "./endpoints";
import { apiInstance } from "./fetchData";
import { setToken } from "./auth";

export async function login(email, password, navigate) {
  apiInstance()
    .post(endpoints.login, {
      email: email,
      password: password,
    })
    .then((response) => {
      setToken(response.data.accessToken);

      location.href=`/`
    })
    .catch((error) => {
      if (error.response) {
        message.error("Invalid Username or Password");
      } else if (error.request) {
        message.error("No Internet Connection");
      } 
    });

}
