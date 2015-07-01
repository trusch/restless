class ServerApp {
    constructor(){
        this.__models = {};
        this.__pages = [];
        this.__debugLogs = true;
        this._initDebugLogs();
    }

    _initDebugLogs(){
        let orig = console.log;
        if(this.__debugLogs){
            console.debug = (...args) => console.log('ServerJS>',...args);
        }else{
            console.debug = (...args) => {};
        }
    }

    Model(name, ...mixins){
        console.debug(`declare server model ${name}...`);
        this.__models[name] = aggregation(ServerModel, ...mixins);
    }

    CreateModel(name, ...args){
        var model = new this.__models[name](...args);
        model.__name = name;
        return model;
    }

    Page(url,renderer){
        console.debug(`declare page for url: ${url}`);
        this.__pages.push({url,renderer});
    }

    renderPage(){
        var renderer = null;
        for(let i=0;i<this.__pages.length;i++){
            if(URL.match(this.__pages[i].url)){
                renderer = this.__pages[i].renderer;
                break;
            }
        }
        if(renderer === null){
            console.debug("no renderer found... sending 404");
            return {
                body: '',
                status: 404,
                header: {}
            };    
        }
        var result = renderer();
        result.body = result.body || '';
        result.status = result.status || 200;
        result.header = result.header || {};
        return result;
    }

}

var app  = new ServerApp();