import React, {useState, useEffect} from 'react';
import ApiService from "services/Network/api";
import {useParams} from 'react-router-dom';
import {Button, Form, Input, message, Select, Spin} from 'antd';

const { Option } = Select;
const URL = "rules"

const layout = {
    labelCol: {span: 2},
    wrapperCol: {span: 8},
};

const DomainEdit = props => {
    let {Uuid} = useParams();
    const [form] = Form.useForm();
    const [roleState, setRoleState] = useState({data: [], loading: false, value: []})
    const [domainState, setDomainState] = useState({data: [], loading: false, value: []})

    const getDomain = value => {
        setDomainState({data: [], loading: true, value: [] });
        ApiService.getAll("domains", {query: value}).then(
            (result) => {
                const data = result.data.domains.map(domain => ({
                    text: domain.name,
                    value: domain.name,
                }));
                setDomainState({data: data, loading: false, value: [] });
            },
            (error) => {
                setDomainState({data: [], loading: false, value: [] });
            }
        )
    };

    const getRole = value => {
        setRoleState({data: [], loading: true, value: [] });
        ApiService.getAll("roles", {query: value}).then(
            (result) => {
                const data = result.data.roles.map(role => ({
                    text: role.title,
                    value: role.title,
                }));
                setRoleState({data: data, loading: false, value: [] });
            },
            (error) => {
                setRoleState({data: [], loading: false, value: [] });
            }
        )
    };

    const initData = {
        domain: "",
        role: "",
        resource: "",
        object: "",
        action: "GET",
        effect: "ALLOW",
    }

    useEffect(() => {
        if (Uuid === undefined) {
            return
        }
        ApiService.get(URL, Uuid).then(
            (result) => {
                form.setFieldsValue({...result.data});
            },
            (error) => {
            }
        )
    }, [form])

    const onFinish = values => {

        const data = {
            ...values,
            domain: values.domain.value,
            role: values.role.value,
        }

        if (Uuid === undefined) {
            // we will create new
            ApiService.post(URL, data).then(
                (result) => {
                    message.info("new " + URL + " created")
                    form.setFieldsValue(initData)
                },
                (error) => {
                    message.error("operation failed ," + error)
                }
            )
        }else{
            // we will update
            ApiService.put(URL,Uuid, data).then(
                (result) => {
                    message.info(" item updated")
                },
                (error) => {
                    message.error("operation failed ," + error)
                }
            )
        }
    };

    const handleRoleChange = value => {
        form.setFieldsValue({role: value})
        setDomainState({data: [], loading: false, value: [] });
    }

    const handleDomainChange = value => {
        form.setFieldsValue({domain: value})
        setDomainState({data: [], loading: false, value: [] });
    }

    return (
        <div className={"section"}>
            <Form {...layout}
                  form={form}
                  name="nest-messages"
                  initialValues={initData}
                  onFinish={onFinish}>
                <Form.Item name="role" label="Role" rules={[{required: true}]}>
                    <Select
                        labelInValue
                        value={roleState.value}
                        placeholder="Select role"
                        notFoundContent={roleState.loading ? <Spin size="small" /> : null}
                        showSearch={true}
                        filterOption={false}
                        onSearch={getRole}
                        onChange={handleRoleChange}
                        style={{ width: '100%' }}
                    >
                        {roleState.data.map(d => (
                            <Option key={d.value}>{d.text}</Option>
                        ))}
                    </Select>
                </Form.Item>
                <Form.Item name="domain" label="Domain" rules={[{required: true}]}>
                    <Select
                        labelInValue
                        value={domainState.value}
                        placeholder="Select domain"
                        notFoundContent={domainState.loading ? <Spin size="small" /> : null}
                        showSearch={true}
                        filterOption={false}
                        onSearch={getDomain}
                        onChange={handleDomainChange}
                        style={{ width: '100%' }}
                    >
                        {domainState.data.map(d => (
                            <Option key={d.value}>{d.text}</Option>
                        ))}
                    </Select>
                </Form.Item>
                <Form.Item name="resource" label="Resource">
                    <Input/>
                </Form.Item>
                <Form.Item name="object" label="Object">
                    <Input/>
                </Form.Item>
                <Form.Item name="action" label="Action">
                    <Select placeholder="Select action">
                        <Option value="OPTION">OPTION</Option>
                        <Option value="GET">GET</Option>
                        <Option value="POST">POST</Option>
                        <Option value="PUT">PUT</Option>
                        <Option value="PATCH">PATCH</Option>
                        <Option value="DELETE">DELETE</Option>
                    </Select>
                </Form.Item>
                <Form.Item name="effect" label="Effect">
                    <Select placeholder="Select effect">
                        <Option value="ALLOW">ALLOW</Option>
                        <Option value="DENY">DENY</Option>
                    </Select>
                </Form.Item>
                <Form.Item wrapperCol={{...layout.wrapperCol, offset: 2}}>
                    <Button type="primary" htmlType="submit">
                        Submit
                    </Button>
                </Form.Item>
            </Form>
        </div>
    )
}

export default DomainEdit;