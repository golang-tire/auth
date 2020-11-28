import React, {useState, useEffect} from 'react';
import ApiService from "services/Network/api";
import {useParams} from 'react-router-dom';
import {Button, Form, Input, InputNumber, Radio, Switch} from 'antd';
import {configs} from 'services/Network/config';

const layout = {
    labelCol: {span: 2},
    wrapperCol: {span: 8},
};

const validateMessages = {
    required: '${label} is required!',
    types: {
        email: '${label} is not a valid email!',
    }
};


const UserEdit = props => {

    let {userUuid} = useParams();
    const [form] = Form.useForm();
    const USERS_URL = configs.API_URL + "/users";
    const initData = {user: {
        firstname: "",
        lastname: "",
        gender: "",
        avatar_url: "",
        username: "",
        password: "",
        email: "",
        enable: true,
        raw_data: ""
    }}

    const getUser = (userUuid) => {
        return ApiService.get(USERS_URL + "/" + userUuid)
    }

    useEffect(() => {
        if (userUuid === undefined) {
            return
        }
        getUser(userUuid).then(
            (result) => {
                form.setFieldsValue({user:{...result.data}});
            },
            (error) => {
            }
        )
    }, [form])

    const onFinish = values => {
        console.log(values);
    };

    return (
        <div>
            <Form {...layout}
                  form={form}
                  name="nest-messages"
                  initialValues={initData}
                  onFinish={onFinish} validateMessages={validateMessages}>
                <Form.Item name={['user', 'username']} label="Username" rules={[{required: true}]}>
                    <Input/>
                </Form.Item>
                <Form.Item name={['user', 'password']} label="Password"
                           rules={[
                               {
                                   required: true,
                                   message: 'Please input your password!',
                               }
                           ]}
                >
                    <Input.Password/>
                </Form.Item>
                <Form.Item name={['user', 'firstname']} label="Firstname">
                    <Input/>
                </Form.Item>
                <Form.Item name={['user', 'lastname']} label="Lastname">
                    <Input/>
                </Form.Item>
                <Form.Item name={['user', 'email']} label="Email" rules={[
                    {
                        required: true,
                        type: 'email'
                    }
                ]}>
                    <Input/>
                </Form.Item>
                <Form.Item label="Gender" name={['user', 'gender']}>
                    <Radio.Group>
                        <Radio.Button value="male">Male</Radio.Button>
                        <Radio.Button value="female">Female</Radio.Button>
                        <Radio.Button value="other">Other</Radio.Button>
                    </Radio.Group>
                </Form.Item>
                <Form.Item name={['user', 'enable']} label="Enable" valuePropName="checked">
                    <Switch />
                </Form.Item>
                <Form.Item name={['user', 'raw_data']} label="Raw data">
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