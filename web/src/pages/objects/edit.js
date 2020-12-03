import React, {useEffect, useState} from 'react';
import ApiService from "services/Network/api";
import {useParams} from 'react-router-dom';
import {Button, Form, Input, message, Select, Spin} from 'antd';

const { Option } = Select;

const ITEMS = "objects"
const URL = "apps/-/objects"

const layout = {
    labelCol: {span: 2},
    wrapperCol: {span: 8},
};

const ObjectEdit = props => {
    let {Uuid} = useParams();
    const [appState, setAppState] = useState({data: [], loading: false})
    const [form] = Form.useForm();

    const initData = {
        identifier: "",
    }

    const getApp = value => {
        setAppState({...appState, loading: true});
        ApiService.getAll("apps", {query: value}).then(
            (result) => {
                const data = result.data.apps.map(app => ({
                    text: app.name,
                    value: app.uuid,
                }));
                setAppState({...appState, data: data, loading: false})
            },
            (error) => {
                setAppState(prevState => {
                    return {...prevState, data: [], loading: false}
                });
            }
        )
    };

    useEffect(() => {
        if (Uuid === undefined) {
            return
        }
        ApiService.get(URL, Uuid).then(
            (result) => {
                // setAppState({...appState, value: result.data.app.name})
                const data = {
                    identifier: result.data.identifier,
                    uuid: result.data.uuid,
                    app:{
                        value: result.data.app.uuid,
                        label: result.data.app.name
                    }
                }
                form.setFieldsValue({...data});
                const value = {value: result.data.uuid, label: result.data.identifier}
                setAppState({...appState, value });
            },
            (error) => {
            }
        )
    }, [form])

    const onFinish = values => {
        if (Uuid === undefined) {
            const url = URL.replace("-", values.app.value)
            // we will create new
            const data = {
                identifier: values.identifier,
            }
            ApiService.post(url, data).then(
                (result) => {
                    message.info("new " + ITEMS + " created")
                    form.setFieldsValue(initData)
                },
                (error) => {
                    message.error("operation failed ," + error)
                }
            )
        }else{
            // we will update
            const data = {
                identifier: values.identifier,
                app_uuid: values.app.value
            }
            ApiService.put(URL, Uuid, data).then(
                (result) => {
                    message.info(" item updated")
                },
                (error) => {
                    message.error("operation failed ," + error)
                }
            )
        }
    };

    const handleAppChange = value => {
        form.setFieldsValue({app: value})
        setAppState({...appState, loading: false});
    }

    return (
        <div className={"section"}>
            <Form {...layout}
                  form={form}
                  name="nest-messages"
                  initialValues={initData}
                  onFinish={onFinish}>
                <Form.Item name="identifier" label="Identifier" rules={[{required: true}]}>
                    <Input/>
                </Form.Item>
                <Form.Item name="app" label="App" rules={[{required: true}]}>
                    <Select
                        labelInValue
                        placeholder="Select app"
                        notFoundContent={appState.loading ? <Spin size="small"/> : null}
                        showSearch={true}
                        filterOption={false}
                        onSearch={getApp}
                        onChange={handleAppChange}
                        style={{width: '100%'}}
                    >
                        {appState.data.map(d => (
                            <Option key={d.value} value={d.value}>{d.text}</Option>
                        ))}
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

export default ObjectEdit;