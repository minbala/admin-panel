import React, { useEffect, useState } from "react";
import { Breadcrumb, Button, Layout, Menu, theme } from "antd";
import { Outlet, Link, redirect, useNavigate } from "react-router-dom";
import { getToken } from "./util/auth";
import Login from "./components/Login";
import "./index.css";
import { apiWithToken } from "./util/fetchData";
import { endpoints } from "./util/endpoints";
const { Header, Content, Footer } = Layout;
const App = () => {
 const token = window?.localStorage.getItem("AccessToken")

 const handleLogout = ()=>{
  apiWithToken().delete(endpoints.logout).then(res=> window?.localStorage.clear()).then((_=>location.href ="/"))
  
 }

  return (
    <Layout className="layout">
      <Header
        style={{
          display: "flex",
          alignItems: "center",
        }}
      >
        <div className="demo-logo" />
        <Menu
          style={{ width: "100%" }}
          theme="dark"
          mode="horizontal"
          items={[
            {
              key: "0",
              label: <Link to="/">Welcome</Link>,
            },
            {
              key: "1",
              label: <Link to="/users">User List</Link>,
            },
            {
              key: "2",
              label: <Link to="/userlogs">User Logs</Link>,
            },
            {
              key: "3",
              label: <Link to="/create">Create User</Link>,
            }, 
           token  && {
              key: "4",
              label: <Button onClick={_=>handleLogout()}>Log Out</Button>,
            },
          ]}
        />
      </Header>
      <Content
        style={{
          padding: "0 50px",
        }}
      >
       {token? <Outlet /> :<Login/>}
      </Content>
    </Layout>
  );
};
export default App;
