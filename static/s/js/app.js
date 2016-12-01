// --------------------------------------------------------------------------------------------------------------------

var app = new Vue({
  el   : '#app',
  data : {
    id        : null,
    idLocal   : null,
    name      : null,
    title     : '',
    author    : '',
    content   : '',
    err       : null,
    state     : 'editing',
    isEditing : true,
    isLoading : false,
  },
  watch: {
    state : function(newState, oldState) {
      this.isEditing = newState === 'editing'
      this.isLoading = newState === 'loading'
    },
  },
  computed : {
    url : function() {
      if ( this.name ) {
        if ( typeof document.location.origin === 'undefined') {
          document.location.origin = document.location.protocol + '//' + document.location.host;
        }
        return document.location.origin + '/' + this.name
      }
      return null
    },
    saveButtonClass : function() {
      if ( this.state === 'editing' ) {
        return 'button is-success is-medium'
      }
      if ( this.state === 'loading' ) {
        return 'button is-success is-medium is-disabled is-loading'
      }
      return 'button is-success is-medium'
    },
  },
  methods : {
    onNew : function() {
      app.id = null
      app.idLocal = null
      app.name = null
      app.title = ''
      app.author = ''
      app.content = ''
      app.err = null
      app.state = 'editing'
    },
    onLoad : function() {
      app.state = 'loading'
      app.err   = null

      var p = axios.get('/api', {
        params : {
          id : app.idLocal,
        },
      })

      p.then(
        function (resp) {
          var data = resp.data
          if ( data.ok ) {
            var payload = data.payload
            app.id      = payload.id
            app.idLocal = payload.id
            app.name    = payload.name
            app.title   = payload.title
            app.author  = payload.author
            app.content = payload.content
            app.err     = null
          }
          else {
            app.err = data.msg
          }
          app.state = 'editing'
        },
        function(err) {
          app.err = 'Error loading article. Please try again later.'
          app.state = 'editing'
        }
      )
    },
    onSave : function() {
      var method
      var data = {
        title   : app.title,
        author  : app.author,
        content : app.content,
      }

      if ( app.name ) {
        // update
        method = 'post'
        data.id = app.id
        data.name = app.name
      }
      else {
        // create
        method = 'put'
      }

      // set to loading
      app.state = 'loading'
      app.err   = null

      var p = axios[method]('/api', data)
      p.then(
        function (resp) {
          var data = resp.data
          if ( data.ok ) {
            // save both the `id` and the `name`
            app.id      = data.id
            app.idLocal = data.id
            app.name    = data.name
          }
          else {
            app.err = data.msg
          }
          app.state = 'editing'
        },
        function(err) {
          app.err = 'Error saving article. Please try again later.'
          app.state = 'editing'
        }
      )
    },
  },
})

// --------------------------------------------------------------------------------------------------------------------
