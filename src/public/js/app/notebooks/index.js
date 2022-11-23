
window.authOkCallback = getNotebooks;

function getNotebooks(token) {
  axios
    .get(window.apiBaseURL+'/api/notebooks', {
      headers: {Authorization: `Bearer ${token}`}
    })
    .then(response => {
      var html = '';
      var reload = false;

      response.data.forEach((item, idx) => {
        reload |= (item.state != 'ACTIVE');

        var menu = '';
        switch (item.menu) {
          case 't4-01': menu = 'NVIDIA T4 1 基 + Intel 2 vCPU'; break;
          case 't4-02': menu = 'NVIDIA T4 1 基 + Intel 4 vCPU'; break;
          case 'a100-01': menu = 'NVIDIA A100 1 基 + Intel 12 vCPU'; break;
        }
        html += '<div class="card">';
          html += '<div class="card-header" id="head-' + item.runtime + '">';
            html += '<h5 class="mb-0">';
              html += '<button class="btn btn-link" data-toggle="collapse" data-target="#collapse-' + item.runtime + '" aria-expanded="' + ((idx == 0) ? 'true' : 'false' )+ '" aria-controls="collapse-' + item.runtime + '">';
                html += item.runtime;
              html += '</button>';
            html += '</h5>';
          html += '</div>';
          html += '<div id="collapse-' + item.runtime + '" class="collapse' + ((idx == 0) ? ' show' : '' )+ '" aria-labelledby="head-' + item.runtime + '" data-parent="#results">';
            html += '<div class="card-body">';
              html += '種別: '+ menu + '<br>';
              html += '状態: '+ item.state + '<br>';
              if (item.proxyUri) {
                html += '接続先: <a href="https://' + item.proxyUri + '/" target="_blank">https://'+ item.proxyUri + '/</a>';
              } else {
                html += '接続先: -';
              }
              html += '<br>';
              html += '<button class="btn btn-secondary" data-toggle="modal" data-target="#notebook-modal"';
                html += 'data-menu="' + item.runtime + '" style="margin-top: 10px;"><span>停止</span></button>&nbsp;&nbsp;';
              html += '<button class="btn btn-danger" data-toggle="modal" data-target="#notebook-modal"';
                html += 'data-menu="' + item.runtime + '" style="margin-top: 10px;"><span>削除</span></button>';
            html += '</div>';
          html += '</div>';
        html += '</div>';
      });
      $("#results").html(html);
      if (reload) setTimeout(function(){getNotebooks(token);}, 5000);
    })
    .catch(error => console.log(error));
}
