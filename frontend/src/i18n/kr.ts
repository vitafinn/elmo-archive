export default {
    window: {
        title: "Elmo Archive"
    },
    server:{
        cn:'China',
        us:'Darkwinter',
        intl:'Haoplay · Global',
        jp:'Haoplay · Japan',
        kr:'Haoplay · Korea',
        tw:'Haoplay · Asia'
    },
    gacha:{
        unknownType: '발주 #{id}',
        type:{
            9: '의상 발주',
            10: '기타 의상 발주',
            '-6': '재회 발주·인형',
            '-7': '재회 발주·군비',
            1:'일반 발주',
            3:'지정 발주',
            4:'군비 강화',
            5:'시작 발주',
            6: '선택 발주·인형',
            7: '선택 발주·군비',
            8:'미스터리 박스'
        },
        statistic:{
            totalCount: 'Total Counter',
            pityCount: 'Pity Counter',
            rank5Data: '5-star Data',
            rank4Data: '4-star Data',
            rank3Data: '3-star Data',
            rank5Avg: 'Avg Pulls per 5-star',
            upRank5Avg: 'Avg Pulls per Up 5-star',
            nonUpRate: 'NonUp 5-star Rate',
        },
        records:{
            title:'Pull Records'
        }
    },
    sync: {
        button: {
            title: 'Synchronize Records',
            tip: 'Pull records from the server and stops when it matches the latest record in the local database'
        },
        loading: 'Syncing...',
        result: {
            success: {
                title: '{server} {uid} Synchronization Success',
                changed: '{count} new records added for {poolType}',
                unchanged: 'No new records added'
            },
            error: {
                cn: 'Synchronization Error(CN)',
                os: 'Synchronization Error(OS)'
            }
        }
    }
}