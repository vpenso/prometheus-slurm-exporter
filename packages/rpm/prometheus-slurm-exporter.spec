%define debug_package %{nil}
%global shortname prometheus-slurm-exporter
%global goipath  github.com/vpenso/prometheus-slurm-exporter
Version:        0.20
%gometa
%global golicenses      LICENSE
%global godocs          README.md

Name:           %{goname}
Release:        %autorelease
Summary:        Prometheus exporter for SLURM metrics
Group:          Monitoring

License:        GPL 3.0
URL:            %{gourl}

Source0:        %{gosource}
Source1:        prometheus-slurm-exporter.service
Source2:        LICENSE
Source3:        README.md

Requires(pre): shadow-utils

Requires(post): systemd
Requires(preun): systemd
Requires(postun): systemd
%{?systemd_requires}
BuildRequires:  systemd-rpm-macros
BuildRequires:  golang(github.com/prometheus/client_golang/prometheus)
BuildRequires:  golang-github-prometheus-common-devel

%description
A Prometheus exporter for metrics extracted from the Slurm resource scheduling system.

%prep
%goprep
%autosetup -N -T -D -a 0 -n %{shortname}-%{version}

%build
make all GOFLAGS=%{gobuildflags}


%install
install -m 0755 -vd %{buildroot}%{_bindir}
install -m 0755 bin/%{shortname} %{buildroot}%{_bindir}

install -m 0755 -vd %{buildroot}%{_unitdir}/
install -m 0755 -vd  %{buildroot}/%{_sharedstatedir}/prometheus
install -m 644 %{SOURCE1} %{buildroot}/%{_unitdir}/%{shortname}.service

%pre
getent group prometheus >/dev/null || groupadd -r prometheus
getent passwd prometheus >/dev/null || \
    useradd -r -g prometheus -d /var/lib/slurm_exporter -s /sbin/nologin \
    -c "Prometheus exporter user" prometheus
exit 0

%post
systemctl enable --now %{shortname}.service
%systemd_post %{shortname}.service

%preun
%systemd_preun %{shortname}.service

%postun
%systemd_postun_with_restart %{shortname}.service

%files
%license LICENSE
%doc README.md
%{_bindir}/%{shortname}
%{_unitdir}/%{shortname}.service
%attr(755, prometheus, prometheus)/%{_sharedstatedir}/prometheus

%changelog
%autochangelog
