import React, {useEffect, useState} from 'react';
import {Button, Dropdown, Input, Menu, message, Modal, Table} from 'antd';
import ApiService from "services/Network/api";
import {configs} from "services/Network/config";
import {Link, useHistory} from "react-router-dom";
import {DeleteOutlined, DownOutlined, ExclamationCircleOutlined, PlusOutlined} from '@ant-design/icons';
import {DropOption} from 'components'

const {confirm} = Modal;

const URL = "domains";
const PageSize = configs.PAGE_SIZE;

const Domains = props => {
    let history = useHistory();
    const [items, setItems] = useState([])
    const [lastQuery, setLastQuery] = useState("")
    const [isLoading, setIsLoading] = useState(false);
    const [showActions, setShowActions] = useState(false);
    const [pagination, setPagination] = useState({
        current: 1,
        pageSize: PageSize,
        total: 0,
    })

    const columns = [
        {
            title: 'Name',
            dataIndex: 'name',
            render: (text, record) => <Link to={`${URL}/edit/${record.uuid}`}>{text}</Link>,
        },
        {
            title: 'Status',
            dataIndex: 'enable',

            render: text => <span>{text ? 'Enable' : 'Disable'}</span>,
        },
        {
            title: 'Created at',
            dataIndex: 'created_at',
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
                            {key: 'delete', name: "Delete"},
                            {key: 'disable', name: "Disable"},
                        ]}
                    />
                )
            },
        },
    ];

    const handleBulkOperationClick = (e) => {
        message.info('Click on menu item.');
        console.log('click', e);
    }

    const bulkMenu = (
        <Menu onClick={handleBulkOperationClick}>
            <Menu.Item key="remove_items" icon={<DeleteOutlined/>}>
                Remove selected {URL}
            </Menu.Item>
        </Menu>
    );

    const loadItems = (page, query="") => {
        setIsLoading(true);
        ApiService.getAll(URL, {limit:PageSize, offset: PageSize * (page-1), query}).then(
            (result) => {
                setItems(result.data[URL]);
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

    const deleteItem = (record) => {
        confirm({
            icon: <ExclamationCircleOutlined/>,
            content: "Are you sure you want to delete `" + record.name + "` ?",
            onOk() {
                ApiService.delete(URL, record.uuid).then(
                    (result) => {
                        message.info("`" + record.name + "` removed")
                        loadItems(pagination.current);
                    },
                    (error) => {
                        message.error("operation failed ," + error)
                    }
                )
            },
        });
    }

    const handleOperationClick = (record, e) => {
        if (e.key === "delete") {
            deleteItem(record)
        }
    }

    // rowSelection object indicates the need for row selection
    const rowSelection = {
        onChange: (selectedRowKeys, selectedRows) => {
            setShowActions(selectedRowKeys.length > 0);
            console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
        },
    };

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
                {showActions && (
                    <Dropdown overlay={bulkMenu} style={{marginRight: "10px"}}>
                        <Button>
                            With selected items <DownOutlined/>
                        </Button>
                    </Dropdown>
                )}
                <Input
                    placeholder="input search text"
                    allowClear
                    onChange={onChange}
                    style={{width: "300px"}}
                />

                <Button onClick={() => {
                    history.push("/" + URL + "/edit")
                }}
                        type="primary" icon={<PlusOutlined/>}
                        style={{float: "right"}}>
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
                onChange={handleTableChange}
            />
        </div>
    );
};

export default Domains;