import { useEffect, useCallback, useState } from "react";
import { BASE_URL } from "./constant";
import axios from "axios";
import { getToken } from "./auth";
import { endpoints } from "./endpoints";
import { message } from "antd";



export function apiWithToken(version = "v1") {
  const apiWithToken = axios.create({
    baseURL: `${endpoints.base}/${version}`,
    timeout: 40000,
    headers: {
      Authorization: `Bearer ${getToken()}`,
    },
  });

  apiWithToken.interceptors.response.use(
    (response) => response,
    (error) => {
      const response = error.response;

      if (response) {
        if (response.status === 401) {
          window.setTokenExpired(true);
        } else {
  console.log(response);
        }
      }

      return Promise.reject(error);
    }
  );
  return apiWithToken;
}



export function useFetchData(
  endpoint,
  refresh,
  version = "v1",
  limit = null,
  offset = null,detail=false,log=false
) {
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);
  const [totalItems, setTotalItems] = useState(0); // New state for total number of items

  useEffect(() => {
    let isMounted = true;
    (async () => {
      try {
        if (endpoint) {
          setLoading(true);
          let requestEndpoint = endpoint;

          if (limit>=0 && offset>=0) {
            requestEndpoint += `?limit=${limit}&offset=${offset}`;
          }
          const res = await apiWithToken(version).get(requestEndpoint);
          const responseData = await res.data;

          if (detail) {
            setData(responseData);
          }else if(log){
            setData(responseData.userLogs);
            setTotalItems(responseData.total);

          }else{setData(responseData.users);
            setTotalItems(responseData.total);
          }
        }
      } catch (error) {
        message.error(error);
      } finally {
        setLoading(false);
      }
    })();
    return () => {
      isMounted = false;
    };
  }, [endpoint, refresh, limit, offset]);

  return { data, error, loading, totalItems }; // Include totalItems in the returned values
}

export function apiInstance() {
  return axios.create({
    baseURL: endpoints.base + "/v1", 
    timeout: 40000,   });
}

