import React, {useEffect, useState} from 'react';
import {Table, Input} from 'antd';
import ApiService from "services/Network/api";
import {Link} from "react-router-dom";

const URL = "audit-logs"
const PageSize = 20;

const AuditLogs = props => {
    const [items, setItems] = useState([]);
    const [lastQuery, setLastQuery] = useState("")
    const [isLoading, setIsLoading] = useState(false);
    const [pagination, setPagination] = useState({
        current: 1,
        pageSize: PageSize,
        total: 0,
    })

    const columns = [
        {
            title: 'User',
            dataIndex: 'user',
            render: (text, record) => <Link to={`users/edit/${record.user.uuid}`}>{record.user.username}</Link>,
        },
        {
            title: 'Object',
            key: "object",
            dataIndex: 'object',
        },
        {
            title: 'Action',
            dataIndex: 'action',
        },
        {
            title: 'Old value',
            dataIndex: 'old_value',
        },
        {
            title: 'New value',
            dataIndex: 'new_value',
        },
        {
            title: 'Created at',
            dataIndex: 'created_at',
            width: 60,
            key: 'created_at',
            fixed: 'right',
        },
    ];

    const loadItems = (page, query="") => {
        setIsLoading(true);
        ApiService.getAll(URL, {limit:PageSize, offset: PageSize * (page-1), query}).then(
            (result) => {
                setItems(result.data[URL.replace("-", "_")]);
                setIsLoading(false);
                setPagination({
                    pageSize: result.data.limit,
                    current: page,
                    total: result.data.total_count,
                })
            },
            (error) => {
                setIsLoading(false);
            }
        )
    }

    useEffect(() => {
        loadItems(pagination.current);
    }, [])

    const handleTableChange = (pagination, filters, sorter) => {
        loadItems(pagination.current);
    };

    const onChange = ({ target: { value } }) => {
        if (value.length >= 3) {
            setLastQuery(value)
            loadItems(pagination.current, value);
        }

        if (lastQuery.length > 0 && value.length < 3){
            setLastQuery("")
            loadItems(pagination.current);
        }
    }
    return (
        <div className={"section"}>
            <div style={{margin: "10px 0px"}}>
                <Input
                    placeholder="input search text"
                    allowClear
                    onChange={onChange}
                    style={{width: "300px"}}
                />
            </div>
            <Table
                bordered={true}
                loading={isLoading}
                rowKey={record => record.uuid}
                pagination={pagination}
                columns={columns}
                dataSource={items}
                onChange={handleTableChange}
            />
        </div>
    );
};

export default AuditLogs;