import React, {useState} from "react";
import AuthService from "services/Auth/authService";
import {useHistory, useLocation, Link} from "react-router-dom";
import {Form, Input, Button, Alert, Row, Col, Card} from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';

const Login = props => {

    let history = useHistory();
    let location = useLocation();
    let { from } = location.state || { from: { pathname: "/" } };

    const [error, setError] = useState(null);
    const onFinish = values => {
        AuthService.login(values.username, values.password)
            .then((res)=>{
                if (res.status === 200){
                    history.replace(from);
                }else{
                    setError(res)
                }
            })
    };

    const OnErrorClose = () => {
        setError(null)
    }

    return (
        <Row type="flex" justify="center" align="middle" style={{minHeight: '100vh'}}>
            <Col span={4}>

            <div className={"login-form-wrapper"}>

                <div className={"login-logo"}>
                 Auth service
                </div>

            <Form
                name="normal_login"
                className="login-form"
                initialValues={{ remember: true }}
                onFinish={onFinish}
            >

                {error && (
                    <div style={{marginBottom: 10}}>
                    <Alert message="Error"
                           afterClose={OnErrorClose}
                           description={"username or password is not correct."}
                           type="error" showIcon closable/>
                    </div>
                 )}
                <Form.Item
                    name="username"
                    rules={[{ required: true, message: 'Please input your Username!' }]}
                >
                    <Input prefix={<UserOutlined className="site-form-item-icon" />} placeholder="Username" />
                </Form.Item>
                <Form.Item
                    name="password"
                    rules={[{ required: true, message: 'Please input your Password!' }]}
                >
                    <Input
                        prefix={<LockOutlined className="site-form-item-icon" />}
                        type="password"
                        placeholder="Password"
                    />
                </Form.Item>
                <Form.Item>
                    {/*<Form.Item name="remember" valuePropName="checked" noStyle>*/}
                    {/*    <Checkbox>Remember me</Checkbox>*/}
                    {/*</Form.Item>*/}
                    <Link className="login-form-forgot" to="/forgot-password">Forgot Password?</Link>
                  </Form.Item>

                <Form.Item>
                    <Button type="primary" htmlType="submit" className="login-form-button">
                        Log in
                    </Button>
                </Form.Item>

            </Form>
            </div>
            </Col>
        </Row>
    )
}

export default Login;