import React, {useState, useEffect} from 'react';
import ApiService from "services/Network/api";
import {useParams} from 'react-router-dom';
import {Button, Form, Input, message, Radio, Switch} from 'antd';
import {configs} from 'services/Network/config';

const URL = "users"

const layout = {
    labelCol: {span: 2},
    wrapperCol: {span: 8},
};

const UserEdit = props => {

    let {Uuid} = useParams();
    const [form] = Form.useForm();
    const initData = {
        firstname: "",
        lastname: "",
        gender: "",
        avatar_url: "",
        username: "",
        password: "",
        email: "",
        enable: true,
        raw_data: ""
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
        if (Uuid === undefined) {
            // we will create new
            ApiService.post(URL, values).then(
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
            ApiService.put(URL,Uuid, values).then(
                (result) => {
                    message.info(" item updated")
                },
                (error) => {
                    message.error("operation failed ," + error)
                }
            )
        }
    };

    return (
        <div className={"section"}>
            <Form {...layout}
                  form={form}
                  name="nest-messages"
                  initialValues={initData}
                  onFinish={onFinish}>
                <Form.Item name="username" label="Username" rules={[{required: true}]}>
                    <Input/>
                </Form.Item>
                <Form.Item name="password" label="Password"
                           rules={[
                               {
                                   required: true,
                                   message: 'Please input your password!',
                               }
                           ]}
                >
                    <Input.Password/>
                </Form.Item>
                <Form.Item name="firstname" label="Firstname">
                    <Input/>
                </Form.Item>
                <Form.Item name="lastname" label="Lastname">
                    <Input/>
                </Form.Item>
                <Form.Item name="email" label="Email" rules={[
                    {
                        required: true,
                        type: 'email'
                    }
                ]}>
                    <Input/>
                </Form.Item>
                <Form.Item label="Gender" name="gender">
                    <Radio.Group>
                        <Radio.Button value="male">Male</Radio.Button>
                        <Radio.Button value="female">Female</Radio.Button>
                        <Radio.Button value="other">Other</Radio.Button>
                    </Radio.Group>
                </Form.Item>
                <Form.Item name="enable" label="Enable" valuePropName="checked">
                    <Switch />
                </Form.Item>
                <Form.Item name="raw_data" label="Raw data">
                    <Input.TextArea/>
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

export default UserEdit;