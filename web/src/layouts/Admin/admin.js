import React, {useState} from "react";
import PropTypes from 'prop-types';
import {Layout, Breadcrumb} from 'antd';
import {useHistory} from "react-router-dom";
import TopHeader from './components/header'
import SideMenu from './components/sideMenu'

import {
    DashboardOutlined,
    UserOutlined,
} from '@ant-design/icons';
import AuthService from "services/Auth/authService";

const { Content, Footer } = Layout;

const Admin = props => {
    let history = useHistory();
    const {children} = props;
    const [collapsed, setCollapsed] = useState(false)
    const [logoUrl, setLogoUrl] = useState("logo.png")
    const [notifications, setNotifications] = useState([])

    const toggle = () =>{
        const newState = !collapsed;
        setCollapsed(newState);
        if (newState=== true) {
            setLogoUrl("logo-small.png")
        }else{
            setLogoUrl("logo.png")
        }
    }

    const handleClick = (e) => {
        if (e.key === "SignOut") {
            AuthService.logout().then(
                (res) => {
                    history.push("/login")
                }
            )
        }
    };

    const onAllNotificationsRead = e => {
        setNotifications([])
    }

    return (
        <Layout className="base-wrapper">
            <SideMenu collapsed={collapsed} logoUrl={logoUrl}/>
            <Layout className="container" id={"mainLayout"}>
                <TopHeader
                    username={""}
                    collapsed={collapsed}
                    onMenuClick={handleClick}
                    notifications={notifications}
                    onClearNotifications={onAllNotificationsRead}
                    OnToggleClick={toggle}
                />
                <Content className="body">
                    <div className="breadcrumb">
                        <Breadcrumb>
                            <Breadcrumb.Item href="">
                                <DashboardOutlined />
                            </Breadcrumb.Item>
                            <Breadcrumb.Item href="">
                                <UserOutlined />
                                <span>User List</span>
                            </Breadcrumb.Item>
                        </Breadcrumb>
                    </div>
                    <div className="content-wrapper">
                        {children}
                    </div>
                </Content>
                <Footer className="footer">Auth ©2020 Created by Mohsen</Footer>
            </Layout>
        </Layout>
    )
};

Admin.propTypes = {
    children: PropTypes.node
};

export default Admin;
