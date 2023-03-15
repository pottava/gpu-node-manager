
$(document).on('show.bs.modal','#vm-modal', function (e) {
  const target = $(e.relatedTarget);
  $('#vm-menu-id').val(target.data('menu'));
  $('#vm-menu').text(target.closest('.card').find('.card-header button').text());
});

$(document).on("click", "#vm-modal .btn-primary", function() {
  $('#vm-modal').modal('hide');

  const params = {menu: $('#vm-menu-id').val()};
  const option = {headers: {Authorization: `Bearer ${idToken}`}};
  axios
    .post(window.apiBaseURL+'/api/vms', params, option)
    .then(_ => location.href = '/vms')
    .catch(error => console.log(error));
});
