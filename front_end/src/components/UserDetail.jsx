import React, { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Descriptions, Empty , Button, Modal, message} from "antd";
import { useFetchData } from "../util/fetchData";
import { endpoints } from "../util/endpoints";
import UpdateUserForm from "./UpdateUserForm";
import { apiWithToken } from "../util/fetchData";
const UserDetail = () => {
  const { id } = useParams();
  const {data} =useFetchData(`${endpoints.user}/${id}`,null,"v1",null,null,true)
  const navigate = useNavigate()
  const [currentData, setCurrentData] = useState({});
  const [openModal, setOpenModal] = useState(false);
  const handleOk = () => {
    setOpenModal(false);
  };
  const handleCancel = () => {
    setOpenModal(false);
  };
  const handleDataUpdate = async (values) => {
    apiWithToken()
      .put("/user", {
        email: values.email,
        name: values.username,
        password: values.password,
        userId: Number(id),
        userRole: data.userRole,
      })
      .then((res) => {setCurrentData({
        email: values.email,
        name: values.username,
        userId: Number(id),
        userRole: data.userRole,
      })
    res.statusText==="OK"&& setOpenModal(false)
    res.statusText==="OK"&& message.success("Success")
    
    });
  };
  const handleDelete = ()=>{
apiWithToken().delete(endpoints.user+"/"+data.id).then(res=>res.status===204&&navigate("/users"))
  }
useEffect(()=>{
  setCurrentData(data)
},[data])
  return (
    <div>
            <Modal
        open={openModal}
        title={"Update User"}
        onOk={handleOk}
        onCancel={handleCancel}
        footer={null}
      >
        <UpdateUserForm data={currentData} updateData={handleDataUpdate} />
      </Modal>
      
{  currentData?    <Descriptions title="User Info" bordered>
        <Descriptions.Item label="UserName">{currentData.name}</Descriptions.Item>
        <Descriptions.Item label="Email">{currentData.email}</Descriptions.Item>
        <Descriptions.Item label="Role">{currentData.userRole}</Descriptions.Item>
        
      </Descriptions>:<Empty/>}
      <Button onClick={(_) => setOpenModal(true)} style={{marginTop:"10px"}}>Edit</Button>
      <Button onClick={(_) => handleDelete()} color="ff0000">Delete</Button>
    </div>
  );
};

export default UserDetail;
