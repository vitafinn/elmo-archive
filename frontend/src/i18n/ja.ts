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
        unknownType: '調達 #{id}',
        type:{
            9: '衣装調達',
            10: 'その他の衣装調達',
            '-6': '再会調達・人形',
            '-7': '再会調達・武器',
            1:'常駐訪問',
            3:'限定訪問',
            4:'軍備拡張',
            5:'スタートダッシュ訪問',
            6: '選択調達・人形',
            7: '選択調達・武器',
            8:'ミステリーボックス'
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
    },
    theme: {
        change: "テーマを変更",
        custom: "カスタムテーマ",
        from: "開始色",
        to: "終了色",
        preview: "プレビュー",
        popular: "人気グラデーション",
        apply: "適用",
        cancel: "キャンセル"
    }
}