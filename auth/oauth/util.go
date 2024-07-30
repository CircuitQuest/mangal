package oauth

// should only be used on implicit grant due to
// access token being part of the fragment
const resendFragmentAsParam = `
<html>
<head>
</head>
<body onload="resend()">

<script>
  function resend(){
    var hash = window.location.hash.substring(1);
    window.location.href = '?' + hash
  }
</script>

</body>
</html>
`
