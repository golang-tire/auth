import React, {Fragment} from "react";
import {Avatar, Menu, Popover, Badge, List, Layout} from 'antd';
import {BellOutlined, MenuFoldOutlined, MenuUnfoldOutlined, RightOutlined} from "@ant-design/icons";
import moment from "moment";

const TopHeader = (props) => {

    return (
        <Layout.Header className={`header ${props.collapsed? "collapsed-header": ""}`} >
            <div className="toggle-btn"
                 onClick={props.OnToggleClick}
            >
                {props.collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
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
                                dataSource={props.notifications}
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
                            {props.notifications.length ? (
                                <div
                                    onClick={props.onClearNotifications}
                                    className={"clear-btn"}
                                >
                                    <span>Clear notifications</span>
                                </div>
                            ) : null}
                        </div>
                    }
                >
                    <Badge
                        count={props.notifications.length}
                        dot
                        className="icon-btn"
                        offset={[-10, 10]}
                    >
                        <BellOutlined/>
                    </Badge>
                </Popover>
                <Menu key="user" mode="horizontal" onClick={props.onMenuClick}>
                    <Menu.SubMenu title={
                        <Fragment>
                            <span style={{ color: '#999', marginRight: 4 }}>Hi,</span>
                            <span>{props.username}</span>
                            <Avatar style={{ marginLeft: 8 }} src="" />
                        </Fragment>
                    }>
                        <Menu.Item key="Profile">Profile</Menu.Item>
                        <Menu.Divider/>
                        <Menu.Item key="SignOut">Sign out</Menu.Item>
                    </Menu.SubMenu>
                </Menu>
            </div>
        </Layout.Header>
    )
}

export default TopHeader;