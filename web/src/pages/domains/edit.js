import React, {useState, useEffect} from 'react';
import ApiService from "services/Network/api";
import {useParams} from 'react-router-dom';
import {Button, Form, Input, Switch} from 'antd';

const URL = "domains"

const layout = {
    labelCol: {span: 2},
    wrapperCol: {span: 8},
};

const DomainEdit = props => {
    let {Uuid} = useParams();
    const [form] = Form.useForm();

    const initData = {domain: {
            name: "",
            enable: true,
        }}

    useEffect(() => {
        if (Uuid === undefined) {
            return
        }
        ApiService.get(URL, Uuid).then(
            (result) => {
                form.setFieldsValue({domain:{...result.data}});
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
                  onFinish={onFinish}>
                <Form.Item name={['domain', 'name']} label="Name" rules={[{required: true}]}>
                    <Input/>
                </Form.Item>
                <Form.Item name={['domain', 'enable']} label="Enable" valuePropName="checked">
                    <Switch />
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