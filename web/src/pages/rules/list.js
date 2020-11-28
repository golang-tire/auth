import React, {useEffect, useState} from 'react';
import { Table , Input , Menu, Dropdown, Button, message} from 'antd';
import {configs} from 'services/Network/config';
import ApiService from "services/Network/api";
import {Link, useHistory} from "react-router-dom";
import { DownOutlined, DeleteOutlined , PlusOutlined} from '@ant-design/icons';
import { DropOption } from 'components'
const { Search } = Input;

const URL = "rules"

const columns = [
    {
        title: 'Role',
        dataIndex: 'role',
        render: (text, record) => <Link to={`${URL}/edit/${record.uuid}`}>{text}</Link>,
    },
    {
        title: 'Domain',
        dataIndex: 'domain',
    },
    {
        title: 'Object',
        dataIndex: 'object',
    },
    {
        title: 'Resource',
        dataIndex: 'resource',
    },
    {
        title: 'Action',
        dataIndex: 'action',
    },
    {
        title: 'Effect',
        dataIndex: 'effect',
    },
    {
        title: '',
        width: 30,
        key: 'operation',
        fixed: 'right',
        render: (text, record) => {
            return (
                <DropOption
                    onMenuClick={e => handleOperationClick(record, e)}
                    menuOptions={[
                        { key: '1', name: "Delete" },
                        { key: '2', name: "Disable" },
                    ]}
                />
            )
        },
    },
];

const handleOperationClick = (record, e) => {
    message.info("click operation " + e.key)
    console.log(e)
}

const handleMenuClick = (e) =>{
    message.info('Click on menu item.');
    console.log('click', e);
}

const menu = (
    <Menu onClick={handleMenuClick}>
        <Menu.Item key="remove_items" icon={<DeleteOutlined />}>
            Remove selected {URL}
        </Menu.Item>
    </Menu>
);


const Rules = props => {
    let history = useHistory();
    const [items, setItems] = useState([])
    const [isLoading, setIsLoading] = useState(false);
    const [showActions, setShowActions] = useState(false);
    const [pagination, setPagination] = useState({
        current: 1,
        pageSize: 10,
        total: 0,
    })

    // rowSelection object indicates the need for row selection
    const rowSelection = {
        onChange: (selectedRowKeys, selectedRows) => {
            setShowActions(selectedRowKeys.length > 0);
            console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
        },
    };

    useEffect(() => {
        setIsLoading(true);
        ApiService.get(URL).then(
            (result) => {
                setItems(result.data[URL]);
                setIsLoading(false);
                setPagination({
                    pageSize: result.data.limit,
                    current: 1,
                    total: result.data.total_count,
                })
            },
            (error) => {
                setIsLoading(false);
            }
        )
    }, [])

    const onSearch = value => console.log(value);
    return (
        <div>
            <div style={{margin: "10px 0px"}}>
                {showActions && (
                    <Dropdown overlay={menu} style={{ marginRight: "10px"}}>
                        <Button>
                            With selected items <DownOutlined />
                        </Button>
                    </Dropdown>
                )}
                <Search
                    placeholder="input search text"
                    allowClear
                    onSearch={onSearch}
                    style={{ width: "300px"}}
                    enterButton
                />

                <Button onClick={()=>{ history.push("/" + URL + "/edit")}}
                        type="primary" icon={<PlusOutlined />}
                        style={{float:"right"}}>
                    Create
                </Button>
            </div>
            <Table
                bordered={true}
                rowSelection={rowSelection}
                loading={isLoading}
                rowKey={record => record.uuid}
                pagination={pagination}
                columns={columns}
                dataSource={items}
            />
        </div>
    );
};

export default Rules;