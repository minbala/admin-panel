import { Button, Form, Input } from "antd";
import React from "react";
import { MailOutlined, LockOutlined, UserOutlined } from "@ant-design/icons";

const UpdateUserForm = ({ data, updateData }) => {
  const [form] = Form.useForm();

  return (
    <Form
      form={form}
      initialValues={
        data && {
          remember: true,
          username: data.name,
          email: data.email,
          password: "",
        }
      }
      onFinish={updateData}
    >
      <Form.Item
        name="username"
        rules={[
          {
            required: true,
            message: "Please type a name",
          },
        ]}
      >
        <Input prefix={<UserOutlined />} placeholder="Name" />
      </Form.Item>

      <Form.Item
        name="email"
        rules={[
          { type: "email", message: "This is not a valid email." },
          { required: true, message: "Please type an email" },
        ]}
      >
        <Input prefix={<MailOutlined />} placeholder="abc@gmail.com" />
      </Form.Item>
      <Form.Item name="password" rules={[]}>
        <Input prefix={<LockOutlined />} placeholder="Password" />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Update User
        </Button>
      </Form.Item>
    </Form>
  );
};

export default UpdateUserForm;
