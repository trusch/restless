restless.Page('^/test$', ()=>{
    
    let apple = restless.CreateModel('Apple');
    apple.getFromUID('27fdd834-d15c-46b6-af0c-606ed0a58d01')

    events.emit('test-visit@websocket')

    return {
        body: `
Url: ${URL}
Method: ${METHOD}
Headers: ${JSON.stringify(HEADERS)}
User: ${USERNAME}
Role: ${ROLE}
Body: ${BODY}

Apple: ${JSON.stringify(apple)}
`
    }
})
