import React, {useState, Fragment} from "react";
import PropTypes from 'prop-types';
import {Layout, Avatar, Menu, Popover, Badge, Breadcrumb, List} from 'antd';
import {NavLink, useHistory} from "react-router-dom";
import {ScrollBar} from 'components'
import moment from 'moment'

import {
    DashboardOutlined,
    BellOutlined,
    MenuUnfoldOutlined,
    MenuFoldOutlined,
    AppstoreAddOutlined,
    UserOutlined,
    TeamOutlined,
    ClusterOutlined,
    AuditOutlined,
    RightOutlined,
} from '@ant-design/icons';
import AuthService from "services/Auth/authService";

const { Header, Content, Footer, Sider } = Layout;
const { SubMenu } = Menu;


const Admin = props => {
    let history = useHistory();
    const {children} = props;
    const [collapsed, setCollapsed] = useState(false)
    const [logoUrl, setLogoUrl] = useState("logo.png")
    const [notifications, setNotifications] = useState([{
        title: "test"
    }])

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
            <Sider trigger={null} collapsible collapsed={collapsed} width={256} className="side-menu">
                <div className="logo">
                    <img src={logoUrl} alt="Auth" height={64} />
                </div>
                <div className="menu-container">
                    <ScrollBar
                        options={{
                            // Disabled horizontal scrolling, https://github.com/utatti/perfect-scrollbar#options
                            suppressScrollX: true,
                        }}
                    >
                        <Menu theme="dark" mode="inline" defaultSelectedKeys={['1']}>
                            <Menu.Item key="1" icon={<DashboardOutlined />}>
                                <NavLink to="/dashboard" activeClassName="active">Dashboard</NavLink>
                            </Menu.Item>
                            <Menu.Item key="2" icon={<AppstoreAddOutlined />}>
                                <NavLink to="/domains" activeClassName="active">Domains</NavLink>
                            </Menu.Item>
                            <Menu.Item key="3" icon={<UserOutlined />}>
                                <NavLink to="/users" activeClassName="active">Users</NavLink>
                            </Menu.Item>
                            <Menu.Item key="4" icon={<TeamOutlined />}>
                                <NavLink to="/roles" activeClassName="active">Roles</NavLink>
                            </Menu.Item>
                            <Menu.Item key="5" icon={<ClusterOutlined />}>
                                <NavLink to="/rules" activeClassName="active" >Rules</NavLink>
                            </Menu.Item>
                            <Menu.Item key="6" icon={<AuditOutlined />}>
                                <NavLink to="/audit-logs" activeClassName="active" >Audit Logs</NavLink>
                            </Menu.Item>
                        </Menu>
                    </ScrollBar>
                </div>
            </Sider>
            <Layout className="container" id={"mainLayout"}>
                <Header className={`header ${collapsed? "collapsed-header": ""}`} >
                    <div className="toggle-btn"
                         onClick={toggle}
                    >
                        {collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
                    </div>
                    <div className="right-container">
                        <Popover
                            placement="bottomRight"
                            trigger="click"
                            key="notifications"
                            getPopupContainer={() => document.querySelector('#mainLayout')}
                            content={
                                <div className={"notifications"}>
                                    <List
                                        itemLayout="horizontal"
                                        dataSource={notifications}
                                        locale={{
                                            emptyText: <span>You have viewed all notifications.</span>,
                                        }}
                                        renderItem={item => (
                                            <List.Item className={"item"}>
                                                <List.Item.Meta
                                                    title={item.title}
                                                    description={moment(item.date).fromNow()}
                                                />
                                                <RightOutlined style={{ fontSize: 10, color: '#ccc' }} />
                                            </List.Item>
                                        )}
                                    />
                                    {notifications.length ? (
                                        <div
                                            onClick={onAllNotificationsRead}
                                            className={"clear-btn"}
                                        >
                                            <span>Clear notifications</span>
                                        </div>
                                    ) : null}
                                </div>
                            }
                        >
                            <Badge
                                count={notifications.length}
                                dot
                                className="icon-btn"
                                offset={[-10, 10]}
                            >
                                <BellOutlined/>
                            </Badge>
                        </Popover>
                        <Menu key="user" mode="horizontal" onClick={handleClick}>
                            <SubMenu title={
                                <Fragment>
                                    <span style={{ color: '#999', marginRight: 4 }}>Hi,</span>
                                    <span>USERNAME</span>
                                    <Avatar style={{ marginLeft: 8 }} src="" />
                                </Fragment>
                            }>
                                <Menu.Item key="Profile">Profile</Menu.Item>
                                <Menu.Divider/>
                                <Menu.Item key="SignOut">Sign out</Menu.Item>
                            </SubMenu>
                        </Menu>
                    </div>
                </Header>
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
