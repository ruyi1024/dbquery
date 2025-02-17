import React, { useState } from 'react';
import { Tooltip } from 'antd';
import { FullscreenOutlined, FullscreenExitOutlined } from '@ant-design/icons';

export default function FullScreen() {
  const [isFullScreen, setIsFullScreen] = useState(true);
  //调用事件
  const fullScreen = () => {
    let isFullScreen = document.webkitIsFullScreen;
    if (!isFullScreen) {
      requestFullScreen();
    } else {
      exitFullscreen();
    }
    setIsFullScreen(isFullScreen);
  };
  //进入全屏
  const requestFullScreen = () => {
    var de = document.documentElement;
    if (de.requestFullscreen) {
      de.requestFullscreen();
    } else if (de.mozRequestFullScreen) {
      de.mozRequestFullScreen();
    } else if (de.webkitRequestFullScreen) {
      de.webkitRequestFullScreen();
    } else if (de.msRequestFullscreen) {
      de.webkitRequestFullScreen();
    }
  };
  //退出全屏
  const exitFullscreen = () => {
    var de = document;
    if (de.exitFullScreen) {
      de.exitFullScreen();
    } else if (de.mozExitFullScreen) {
      de.mozExitFullScreen();
    } else if (de.webkitExitFullscreen) {
      de.webkitExitFullscreen();
    } else if (de.msExitFullscreen) {
      de.msExitFullscreen();
    }
  };
  return (
    <Tooltip placement="bottom" title={<span>{isFullScreen ? '全屏' : '退出全屏'}</span>}>
           {isFullScreen ? (
        <FullscreenOutlined
          style={{
            fontSize: '20px',
          }}
          onClick={fullScreen}
        />
      ) : (
        <FullscreenExitOutlined
          style={{
            fontSize: '20px',
          }}
          onClick={fullScreen}
        />
      )}
    </Tooltip>
  );
}
