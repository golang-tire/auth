import React from "react";
import {ScrollBar} from "components";
import {Menu, Layout} from "antd";
import {NavLink} from "react-router-dom";
import settings from "settings";
import iconMap from "utils/iconMap";

const SideMenu = (props) => {

    return (
        <Layout.Sider trigger={null} collapsible collapsed={props.collapsed} width={256} className="side-menu">
            <div className="logo">
                <img src={props.logoUrl} alt="Auth" height={64}/>
            </div>
            <div className="menu-container">
                <ScrollBar
                    options={{
                        // Disabled horizontal scrolling, https://github.com/utatti/perfect-scrollbar#options
                        suppressScrollX: true,
                    }}
                >
                    <Menu theme="dark" mode="inline" defaultSelectedKeys={settings.defaultSelectedMenus}>
                        {
                            settings.routeList.map(item => {
                                if (item.sideMenu){
                                    return(
                                        <Menu.Item key={item.id} icon={iconMap[item.icon]}>
                                            <NavLink to={item.path}>{item.name}</NavLink>
                                        </Menu.Item>
                                    )
                                }
                            })
                        }
                    </Menu>
                </ScrollBar>
            </div>
        </Layout.Sider>
    )
}

export default SideMenu;