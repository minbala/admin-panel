import { useNavigate } from "react-router-dom";
import { Button, Card, Input, Space, Table } from "antd";
import { useFetchData } from "../util/fetchData";
import { endpoints } from "../util/endpoints";
import { useRef, useState } from "react";
import { SearchOutlined } from "@ant-design/icons";

export default function UsersList(params) {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const navigate = useNavigate();
 
  const { data, error, loading, totalItems } = useFetchData(
    `${endpoints.user}`,
    null,
    "v1",
    pageSize, // Conditionally set the pageSize based on limitNoffset prop
    (currentPage - 1) * pageSize // Conditionally set the offset based on limitNoffset prop
  );

  const handlePageChange = (page, pageSize) => {
    setCurrentPage(page);
    setPageSize(pageSize);
  };
 

  const columns = [
    {
      title: "Id",
      dataIndex: "id",
      key: "Id",
    },
    {
      title: "Name",
      dataIndex: "name",
      key: "name",
    },
    {
      title: "Email",
      dataIndex: "email",
      key: "email",
    },
    {
      title: "Created At",
      dataIndex: "createdAt",
      key: "createdat",
      render:(date)=>{
        function formatDateToDdmmyy(dateString) {
          const date = new Date(dateString);
        
          const day = String(date.getDate()).padStart(2, "0");
          const month = String(date.getMonth() + 1).padStart(2, "0");
          const year = String(date.getFullYear()).slice(-2);
        
          return `${day}.${month}.${year}`;
        }
        return formatDateToDdmmyy(date)
      }
    },
    
  ];
  const handleOnRow = (e) => {
    return {
      onClick: () => {
        navigate(`/user/${e.id}`);
      },
    };
  };
  if (error) {
    return "error";
  }

  return (
    <Card title="Users">
      
      <div>
        <Table
          loading={loading}
          dataSource={data || []}
          columns={columns}
          rowClassName={"cursor-pointer"}
          onRow={handleOnRow}
          pagination={{
            current: currentPage,
            pageSize: pageSize,
            total: totalItems,
            onChange: handlePageChange,
          }}
        />
      </div>
    </Card>
  );
}
