import React, {useState, useEffect} from 'react';
import ApiService from "services/Network/api";
import {useParams} from 'react-router-dom';
import {Button, Form, Input, Switch} from 'antd';
import {configs} from 'services/Network/config';

const URL = "rules"

const layout = {
    labelCol: {span: 2},
    wrapperCol: {span: 8},
};

const DomainEdit = props => {
    let {Uuid} = useParams();
    const [form] = Form.useForm();

    const initData = {rule: {
            name: "",
            enable: true,
        }}

    const getItem = (Uuid) => {
        return ApiService.get(configs.API_URL + "/" + URL + "/" + Uuid)
    }

    useEffect(() => {
        if (Uuid === undefined) {
            return
        }
        getItem(Uuid).then(
            (result) => {
                form.setFieldsValue({rule:{...result.data}});
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
                <Form.Item name={['rule', 'role']} label="Role" rules={[{required: true}]}>
                    <Input/>
                </Form.Item>
                <Form.Item name={['rule', 'enable']} label="Enable" valuePropName="checked">
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