import React, {useState, Fragment} from "react";
import PropTypes from 'prop-types';
import {Layout, Avatar, Menu, Popover, Badge, List, Breadcrumb} from 'antd';
import {NavLink, useLocation} from "react-router-dom";
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
} from '@ant-design/icons';

const { Header, Content, Footer, Sider } = Layout;
const { SubMenu } = Menu;

const Admin = props => {
    const {children} = props;
    const [collapsed, setCollapsed] = useState(false)
    const location = useLocation();

    const toggle = () =>{
        setCollapsed(!collapsed);

        console.log(location.pathname);
    }

    return (
        <Layout>
            <Sider trigger={null} collapsible collapsed={collapsed} width={256}
                style={{
                    height: '100vh',
                }}
            >
                <div className="logo">
                    Auth
                </div>
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
            </Sider>
            <Layout className="site-layout">
                <Header className="header-background" style={{ padding: 0 }}>
                    <div style={{float:"left"}}>
                        {React.createElement(collapsed ? MenuUnfoldOutlined : MenuFoldOutlined, {
                            className: 'trigger',
                            onClick: toggle,
                        })}
                    </div>
                    <div style={{float:"right"}}>
                        <Popover
                            placement="bottomRight"
                            trigger="click"
                            key="notifications"
                        >
                            <Badge
                                count="10"
                                dot
                                offset={[-10, 10]}
                            >
                                <BellOutlined/>
                            </Badge>
                        </Popover>
                        <Menu key="user" mode="horizontal" style={{float:"right"}}>
                            <SubMenu title={
                                <Fragment>
                                    <span style={{ color: '#999', marginRight: 4 }}>Hi,</span>
                                    <span>Mohsen</span>
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
                <Breadcrumb style={{margin: 24}}>
                    <Breadcrumb.Item href="">
                        <DashboardOutlined />
                    </Breadcrumb.Item>
                    <Breadcrumb.Item href="">
                        <UserOutlined />
                        <span>User List</span>
                    </Breadcrumb.Item>
                </Breadcrumb>
                <Content style={{ margin: '0 16px 16px 16px', padding: '0 12px 12px 0', minHeight: 280}}>
                        {children}
                </Content>
                <Footer style={{ textAlign: 'center' }}>Auth ©2020 Created by Mohsen</Footer>
            </Layout>
        </Layout>
    )
};

Admin.propTypes = {
    children: PropTypes.node
};

export default Admin;
