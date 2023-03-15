
window.authOkCallback = getVMs;

function getVMs (token) {
  axios
    .get(window.apiBaseURL+'/api/vms', {
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
          case 'cpu-01': menu = 'Intel 2 vCPU'; break;
          case 't4-01': menu = 'NVIDIA T4 1 基 + Intel 2 vCPU'; break;
        }
        html += '<div class="card">';
          html += '<div class="card-header" id="head-' + item.name + '">';
            html += '<h5 class="mb-0">';
              html += '<button class="btn btn-link" data-toggle="collapse" data-target="#collapse-' + item.name + '" aria-expanded="' + ((idx == 0) ? 'true' : 'false' )+ '" aria-controls="collapse-' + item.name + '">';
                html += item.name;
              html += '</button>';
            html += '</h5>';
          html += '</div>';
          html += '<div id="collapse-' + item.name + '" class="collapse" aria-labelledby="head-' + item.name + '" data-parent="#results">';
            html += '<div class="card-body">';
              html += '種別: '+ menu + '<br>';
              html += '<br>作成: '+ item.created_at + '<br>';

              switch (item.state) {
                case 'true':
                  html += '<br>';
                  html += '<button class="btn btn-danger" data-toggle="modal" data-target="#vm-modal"';
                    html += 'data-type="delete" data-name="' + item.name + '" style="margin-top: 10px;">';
                    html += '<span>削除</span></button>';
                  break;
                case 'false':
                  html += '削除: '+ item.updated_at;
                  break;
              }
            html += '</div>';
          html += '</div>';
        html += '</div>';
      });
      $("#results").html(html);
    })
    .catch(error => console.log(error));
}

$(document).on('show.bs.modal','#vm-modal', function (e) {
  const target = $(e.relatedTarget);
  $('#request-type').val(target.data('type'));
  $('#vm-id').val(target.data('name'));
  $('#vm h4').text(target.closest('.card').find('.card-header button').text());
  $('#vm span').text(target.find('span').text());
});

$(document).on("click", "#vm-modal .btn-primary", function() {
  $('#vm-modal').modal('hide');

  const params = {id: $('#vm-id').val()};
  const option = {Authorization: `Bearer ${idToken}`};
  switch ($('#request-type').val()) {
  case 'delete':
    axios
      .delete(window.apiBaseURL+'/api/vms', {headers: option, data: params})
      .then(_ => location.href = '/vms')
      .catch(error => console.log(error));
    break;
  default:
    params["action"] = $('#request-type').val();
    axios
      .put(window.apiBaseURL+'/api/vms', params, {headers: option})
      .then(_ => location.href = '/vms')
      .catch(error => console.log(error));
  }
});
