import extra from './zh-CN-extra'
export default {
    ...extra,
    window: {
        title: "Elmo Archive"
    },
    server: {
        cn: '国服',
        us:'Darkwinter',
        intl:'Haoplay · 全球服',
        jp:'Haoplay · 日服',
        kr:'Haoplay · 韩服',
        tw:'Haoplay · 亚服'
    },
    gacha: {
        unknownType: '采购 #{id}',
        type: {
            9: '外观采购',
            10: '其他外观采购',
            '-6': '重逢采购·人形',
            '-7': '重逢采购·武器',
            1: '常规采购',
            3: '定向采购',
            4: '军备提升',
            5: '初始采购',
            6: '自选采购·人形',
            7: '自选采购·军备',
            8: '神秘箱'
        },
        statistic: {
            totalCount: '记录总数',
            pityCount: '保底进度',
            rank5Data: '五星数据',
            rank4Data: '四星数据',
            rank3Data: '三星数据',
            rank5Avg: '五星平均抽数',
            upRank5Avg: 'Up五星平均抽数',
            nonUpRate: '五星歪率',
        },
        records: {
            title: '抽卡记录',
            tip: '物品名称来自已安装的游戏数据，可能与界面语言不同。没有简体中文名称时，软件会使用繁体中文。'
        }
    },
    sync: {
        button: {
            title: '同步记录',
            tip: '从服务器拉取抽卡记录，匹配到本地数据库最新记录后停止'
        },
        loading: '正在同步记录...',
        result: {
            success: {
                title: '{server} {uid} 同步成功',
                changed: '{poolType} 新增 {count} 条',
                unchanged: '无新增数据'
            },
            error: {
                cn: '国服同步出错',
                os: '国际服同步出错'
            }
        }
    },
    version: {
        update: {
            notify: '有新的版本可用，是否更新？',
            latest: '已是最新版本',
            confirm: '是',
            cancel: '否'
        },
    },
    theme: {
        change: "更改主题",
        custom: "自定义主题",
        from: "起始颜色",
        to: "结束颜色",
        preview: "预览效果",
        popular: "流行渐变",
        apply: "应用",
        cancel: "取消"
    }
}
