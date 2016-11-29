// --------------------------------------------------------------------------------------------------------------------

var app = new Vue({
  el   : '#app',
  data : {
    id      : null,
    name    : null,
    title   : '',
    author  : '',
    content : '',
    err     : null,
    state   : 'editing',
  },
  computed: {
    url : function() {
      if ( this.name ) {
        return 'https://publish.li/' + this.name
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
  methods: {
    onLoad: function() {
      app.state = 'loading'
      app.err   = null

      var p = axios.get('/post', {
        id : app.id,
      })

      p.then(
        function (resp) {
          console.log(resp)
          var data = resp.data
          if ( data.ok ) {
            app.id      = data.id
            app.name    = data.name
            app.title   = data.title
            app.author  = data.author
            app.content = data.content
            app.err     = null
          }
          else {
            app.err = data.msg
          }
          app.state = 'editing'
        },
        function(err) {
          console.log(err)
          app.err = 'Error saving article. Please try again later.'
          app.state = 'editing'
        }
      )
    },
    onSave: function() {
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
          console.log(resp)
          var data = resp.data
          if ( data.ok ) {
            // save both the `id` and the `name`
            app.id   = data.id
            app.name = data.name
          }
          else {
            app.err = data.msg
          }
          app.state = 'editing'
        },
        function(err) {
          console.log(err)
          app.err = 'Error saving article. Please try again later.'
          app.state = 'editing'
        }
      )
    },
  },
})

// --------------------------------------------------------------------------------------------------------------------
