layui.define(['laytpl', 'layer'], function (exports) {
  var $ = layui.jquery
    , laytpl = layui.laytpl
    , layer = layui.layer
    , setter = layui.setter
    , device = layui.device()
    , hint = layui.hint()
    , view = function (id) {
    return new Class(id)
  }
    , SHOW = 'layui-show'
    , LAY_BODY = 'LAY_app_body'
    , Class = function (id) {
    this.id = id
    this.container = $('#' + (id || LAY_BODY))
  }

  view.loading = function (elem) {
    elem.append(
      this.elemLoad = $(
        '<i class="layui-anim layui-anim-rotate layui-anim-loop layui-icon layui-icon-loading layadmin-loading"></i>'),
    )
  }

  view.removeLoad = function () {
    this.elemLoad && this.elemLoad.remove()
  }

  view.exit = function (callback) {

    layui.data(setter.tableName, {
      key: setter.request.tokenName
      , remove: true,
    })

    callback && callback()
  }

  view.req = function (options) {
    var that = this
      , success = options.success
      , error = options.error
      , request = setter.request
      , response = setter.response
      , debug = function () {
      return setter.debug
        ? '<br><cite>URL：</cite>' + options.url
        : ''
    }

    options.data = options.data || {}
    options.headers = options.headers || {}
``
    if (request.tokenName) {
      var sendData = typeof options.data === 'string'
        ? JSON.parse(options.data)
        : options.data

      options.data[request.tokenName] = request.tokenName in sendData
        ? options.data[request.tokenName]
        : (layui.data(setter.tableName)[request.tokenName] || '')

      options.headers[request.tokenName] = request.tokenName in
      options.headers
        ? options.headers[request.tokenName]
        : (layui.data(setter.tableName)[request.tokenName] || '')
    }

    delete options.success
    delete options.error

    return $.ajax($.extend({
      type: 'get'
      , dataType: 'json'
      , success: function (res) {
        var statusCode = response.statusCode

        if (res[response.statusName] == statusCode.ok) {
          typeof options.done === 'function' && options.done(res)
        } else if (res[response.statusName] == statusCode.logout) {
          view.exit()
        } else {
          const error = [
            '<cite>Error：</cite> ' + (res[response.msgName] || '返回状态码异常'),
          ].join('')
          layer.msg(error)
          return false
        }

        typeof success === 'function' && success(res)
      }
      , error: function (e, code) {
        var error = [
          '请求异常，请重试<br><cite>错误信息：</cite>' + code
          , debug(),
        ].join('')
        layer.msg(error)

        typeof error === 'function' && error(res)
      },
    }, options))
  }

  view.popup = function (options) {
    var success = options.success
      , skin = options.skin

    delete options.success
    delete options.skin

    return layer.open($.extend({
      type: 1
      , title: '提示'
      , content: ''
      , id: 'LAY-system-view-popup'
      , shadeClose: true
      , closeBtn: false
      , success: function (layero, index) {
        var elemClose = $('<i class="layui-icon" close>&#x1006;</i>')
        layero.append(elemClose)
        elemClose.on('click', function () {
          layer.close(index)
        })
        typeof success === 'function' && success.apply(this, arguments)
      },
    }, options))
  }

  view.error = function (content, options) {
    return view.popup($.extend({
      content: content
      , maxWidth: 300
      , offset: 't'
      , anim: 6
      , id: 'LAY_adminError',
    }, options))
  }

  Class.prototype.render = function (views, params) {
    var that = this, router = layui.router()
    views = setter.views + views + setter.engine

    $('#' + LAY_BODY).children('.layadmin-loading').remove()
    view.loading(that.container)

    $.ajax({
      url: views
      , type: 'get'
      , dataType: 'html'
      , data: {
        v: layui.cache.version,
      }
      , success: function (html) {
        html = '<div>' + html + '</div>'

        var elemTitle = $(html).find('title')
          , title = elemTitle.text() ||
          (html.match(/\<title\>([\s\S]*)\<\/title>/) || [])[1]

        var res = {
          title: title
          , body: html,
        }

        elemTitle.remove()
        that.params = params || {}

        if (that.then) {
          that.then(res)
          delete that.then
        }

        that.parse(html)
        view.removeLoad()

        if (that.done) {
          that.done(res)
          delete that.done
        }

      }
      , error: function (e) {
        view.removeLoad()

        if (that.render.isError) {
          return view.error('请求视图文件异常，状态：' + e.status)
        }

        if (e.status === 404) {
          that.render('template/tips/404')
        } else {
          that.render('template/tips/error')
        }

        that.render.isError = true
      },
    })
    return that
  }

  Class.prototype.parse = function (html, refresh, callback) {
    var that = this
      , isScriptTpl = typeof html === 'object'
      , elem = isScriptTpl ? html : $(html)
      , elemTemp = isScriptTpl ? html : elem.find('*[template]')
      , fn = function (options) {
      var tpl = laytpl(options.dataElem.html())
        , res = $.extend({
        params: router.params,
      }, options.res)

      options.dataElem.after(tpl.render(res))
      typeof callback === 'function' && callback()

      try {
        options.done && new Function('d', options.done)(res)
      } catch (e) {
        console.error(options.dataElem[0], '\n存在错误回调脚本\n\n', e)
      }
    }
      , router = layui.router()

    elem.find('title').remove()
    that.container[refresh ? 'after' : 'html'](elem.children())

    router.params = that.params || {}

    for (var i = elemTemp.length; i > 0; i--) {
      (function () {
        var dataElem = elemTemp.eq(i - 1)
          , layDone = dataElem.attr('lay-done') ||
          dataElem.attr('lay-then')
          , url = laytpl(dataElem.attr('lay-url') || '').
          render(router)
          , data = laytpl(dataElem.attr('lay-data') || '').
          render(router)
          , headers = laytpl(dataElem.attr('lay-headers') || '').
          render(router)

        try {
          data = new Function('return ' + data + ';')()
        } catch (e) {
          hint.error('lay-data: ' + e.message)
          data = {}
        }

        try {
          headers = new Function('return ' + headers + ';')()
        } catch (e) {
          hint.error('lay-headers: ' + e.message)
          headers = headers || {}
        }

        if (url) {
          view.req({
            type: dataElem.attr('lay-type') || 'get'
            , url: url
            , data: data
            , dataType: 'json'
            , headers: headers
            , success: function (res) {
              fn({
                dataElem: dataElem
                , res: res
                , done: layDone,
              })
            },
          })
        } else {
          fn({
            dataElem: dataElem
            , done: layDone,
          })
        }
      }())
    }

    return that
  }

  Class.prototype.autoRender = function (id, callback) {
    var that = this
    $(id || 'body').
      find('*[template]').
      each(function (index, item) {
        var othis = $(this)
        that.container = othis
        that.parse(othis, 'refresh')
      })
  }

  Class.prototype.send = function (views, data) {
    var tpl = laytpl(views || this.container.html()).
      render(data || {})
    this.container.html(tpl)
    return this
  }

  Class.prototype.refresh = function (callback) {
    var that = this
      , next = that.container.next()
      , templateid = next.attr('lay-templateid')

    if (that.id != templateid) return that

    that.parse(that.container, 'refresh', function () {
      that.container.siblings(
        '[lay-templateid="' + that.id + '"]:last').
        remove()
      typeof callback === 'function' && callback()
    })

    return that
  }

  Class.prototype.then = function (callback) {
    this.then = callback
    return this
  }
  Class.prototype.done = function (callback) {
    this.done = callback
    return this
  }
  exports('view', view)
})