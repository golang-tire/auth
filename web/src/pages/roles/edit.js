import React, {useState, useEffect} from 'react';
import ApiService from "services/Network/api";
import {useParams} from 'react-router-dom';
import {Button, Form, Input, message, Switch} from 'antd';

const URL = "roles"

const layout = {
    labelCol: {span: 2},
    wrapperCol: {span: 8},
};

const DomainEdit = props => {
    let {Uuid} = useParams();
    const [form] = Form.useForm();

    const initData = {
        title: "",
        enable: true,
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
                <Form.Item name="title" label="Name" rules={[{required: true}]}>
                    <Input/>
                </Form.Item>
                <Form.Item name="enable" label="Enable" valuePropName="checked">
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