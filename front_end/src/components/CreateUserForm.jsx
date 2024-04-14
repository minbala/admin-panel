import { Button, Card, Form, Input, Select, message } from "antd";
import React from "react";
import { apiWithToken } from "../util/fetchData";
import { endpoints } from "../util/endpoints";

const CreateUserForm = () => {
  const [form] = Form.useForm();
  const handleCreateUser = async(values) => {
apiWithToken().post(endpoints.user,values).then(res=>{

    
    if(res.status===201){ message.success("Success")

form.resetFields()}else{
        message.error("failed")
    }

})
  };
  return (
    <Card title={"Create User"}>
      <Form form={form} onFinish={handleCreateUser}>
        <Form.Item
          name="name"
          rules={[{ required: true, message: "Name is required" }]}
        >
          <Input placeholder="Please input name" />
        </Form.Item>
        <Form.Item
          name="email"
          rules={[
            { required: true, message: "Email is required" },
            { type: "email", message: "Email is invalid" },
          ]}
        >
          <Input placeholder="Please input email" />
        </Form.Item>
        <Form.Item
          name="password"
          rules={[{ required: true, message: "Password is required" }]}
        >
          <Input placeholder="Please input Password" />
        </Form.Item>
        <Form.Item
          name="userRole"
          rules={[{ required: true, message: "Role is required" }]}
        >
          <Select
          placeholder="Please providerole."
            style={{
              width: 160,
            }}
            options={[
              {
                value: "normal user",
                label: "User",
              },
              {
                value: "admin",
                label: "Admin",
              },
            ]}
          />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit">
            Create User
          </Button>
        </Form.Item>
      </Form>
    </Card>
  );
};

export default CreateUserForm;
