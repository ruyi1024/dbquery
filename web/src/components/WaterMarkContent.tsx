import React from 'react'
import { ReactNode, useMemo } from 'react'

type svgPropsType = {
    text?: string
    fontSize?: number
    fillOpacity?: number
    fillColor?: string
}
const SvgTextBg = (props: svgPropsType) => {
    const { text = 'waterMark', fontSize = 14, fillOpacity = '0.1', fillColor = '#000' } = props
    const res = `
      <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="180px" height="180px" viewBox="0 0 180 180">
        <text x="-100" y="-30" fill='${fillColor}'  transform = "rotate(-35 220 -220)" fill-opacity='${fillOpacity}' font-size='${fontSize}'> ${text}</text>
      </svg>`

    const blob = new Blob([res], {
        type: 'image/svg+xml',
    })

    const url = URL.createObjectURL(blob)

    return (
        <div
            style={{
                position: 'absolute',
                width: '100%',
                height: '100%',
                backgroundImage: `url(${url})`,
                top: 0,
                left: 0,
                zIndex: 999,
                pointerEvents: 'none', //点击穿透
            }}
        ></div>
    )
}

type propsType = {
    children?: ReactNode
} & Partial<svgPropsType>

const WaterMarkContent = (props: propsType) => {
    const { text, fontSize, fillOpacity, fillColor } = props

    const memoInfo = useMemo(
        () => ({
            text,
            fontSize,
            fillOpacity,
            fillColor,
        }),
        [text, fontSize, fillOpacity, fillColor]
    )
    return (
        <div style={{ position: 'relative', width: '100%', height: ' 100%' }}>
            {props.children}
            <SvgTextBg {...memoInfo} />
        </div>
    )
}

export default WaterMarkContent

