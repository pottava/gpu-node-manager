
window.authOkCallback = getNotebooks;

function getNotebooks (token) {
  axios
    .get(window.apiBaseURL+'/api/notebooks', {
      headers: {Authorization: `Bearer ${token}`}
    })
    .then(response => {
      let html = '';
      let reload = false;

      response.data.sort((a, b) => {
        if (a.created_at < b.created_at) {
          return 1;
        }
        if (a.created_at > b.created_at) {
          return -1;
        }
        return 0;

      }).forEach((item, idx) => {
        reload |= (item.state != 'ACTIVE') && (item.state != 'DELETED');

        let menu = '';
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
              if (item.state == 'ACTIVE' && item.proxyUri) {
                html += '接続先: <a href="https://' + item.proxyUri + '/" target="_blank">https://'+ item.proxyUri + '/</a>';
              } else {
                html += '接続先: -';
              }
              html += '<br>作成: '+ item.created_at + '<br>';

              switch (item.state) {
              case 'ACTIVE':
                html += '<br>';
                html += '<button class="btn btn-secondary" data-toggle="modal" data-target="#notebook-modal"';
                  html += 'data-type="stop" data-menu="' + item.runtime + '" style="margin-top: 10px;">';
                  html += '<span>停止</span></button>&nbsp;&nbsp;';
                html += '<button class="btn btn-danger" data-toggle="modal" data-target="#notebook-modal"';
                  html += 'data-type="delete" data-menu="' + item.runtime + '" style="margin-top: 10px;">';
                  html += '<span>削除</span></button>';
                break;
              case 'STOPPED':
                html += '<br>';
                html += '<button class="btn btn-secondary" data-toggle="modal" data-target="#notebook-modal"';
                  html += 'data-type="start" data-menu="' + item.runtime + '" style="margin-top: 10px;">';
                  html += '<span>再開</span></button>&nbsp;&nbsp;';
                html += '<button class="btn btn-danger" data-toggle="modal" data-target="#notebook-modal"';
                  html += 'data-type="delete" data-menu="' + item.runtime + '" style="margin-top: 10px;">';
                  html += '<span>削除</span></button>';
                break;
              }
            html += '</div>';
          html += '</div>';
        html += '</div>';
      });
      $("#results").html(html);
      if (reload) setTimeout(function(){getNotebooks(token);}, 10*1000);
    })
    .catch(error => console.log(error));
}

$(document).on('show.bs.modal','#notebook-modal', function (e) {
  const target = $(e.relatedTarget);
  $('#request-type').val(target.data('type'));
  $('#notebook-id').val(target.data('menu'));
  $('#notebook h4').text(target.closest('.card').find('.card-header button').text());
  $('#notebook span').text(target.find('span').text());
});

$(document).on("click", "#notebook-modal .btn-primary", function() {
  $('#notebook-modal').modal('hide');

  const params = {id: $('#notebook-id').val()};
  const option = {Authorization: `Bearer ${idToken}`};
  switch ($('#request-type').val()) {
  case 'delete':
    axios
      .delete(window.apiBaseURL+'/api/notebooks', {headers: option, data: params})
      .then(_ => location.href = '/notebooks')
      .catch(error => console.log(error));
    break;
  default:
    params["action"] = $('#request-type').val();
    axios
      .put(window.apiBaseURL+'/api/notebooks', params, {headers: option})
      .then(_ => location.href = '/notebooks')
      .catch(error => console.log(error));
  }
});
