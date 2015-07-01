app.Page('^/$', ()=>{
    return {
        status: 301,
        header: {
            location: '/index.js'
        }
    };
});