import React, {useState} from "react";
import {Avatar, Badge, Layout, Menu, Popover} from "antd";
import {BellOutlined, MenuFoldOutlined, MenuUnfoldOutlined} from "@ant-design/icons";

const Header = props => {
    const [username] = props
    const [collapsed, setCollapsed] = useState(props.collapsed);

    const toggle = () =>{
        setCollapsed(!collapsed);
    }

    return (
        <Layout.Header>
            <div style={{float:"left"}}>
                {React.createElement(collapsed ? MenuUnfoldOutlined : MenuFoldOutlined, {
                    className: 'trigger',
                    onClick: toggle,
                })}
            </div>
            <div style={{float:"right"}}>
                <Menu key="user" mode="horizontal" style={{float:"right"}} onClick={handleClick}>
                    <SubMenu title={
                        <Fragment>
                            <span style={{ color: '#999', marginRight: 4 }}>Hi,</span>
                            <span>{username}</span>
                            <Avatar style={{ marginLeft: 8 }} src="" />
                        </Fragment>
                    }>
                        <Menu.Item key="Profile">Profile</Menu.Item>
                        <Menu.Divider/>
                        <Menu.Item key="SignOut">Sign out</Menu.Item>
                    </SubMenu>
                </Menu>
            </div>
        </Layout.Header>
    )
}