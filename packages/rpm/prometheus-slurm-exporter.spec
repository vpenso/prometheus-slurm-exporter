%define        __spec_install_post %{nil}
%define          debug_package %{nil}
%define        __os_install_post %{_dbpath}/brp-compress

Name:           prometheus-slurm-exporter
Version:        0.3
Release:        1%{?dist}
Summary:        Prometheus exporter for SLURM metrics
Group:          Monitoring

License:        GPL 3.0
URL:            https://github.com/vpenso/prometheus-slurm-exporter

Source0:        https://github.com/vpenso/prometheus-slurm-exporter/releases/download/%{version}/slurm_exporter-%{version}.linux-amd64.tar.gz
Source1:        prometheus-slurm-exporter.service
Source2:        LICENSE
Source3:        README.md

Requires(pre): shadow-utils

Requires(post): systemd
Requires(preun): systemd
Requires(postun): systemd
%{?systemd_requires}
BuildRequires:  systemd

BuildRoot:      %{_tmppath}/%{name}-%{version}-1-root

%description
A Prometheus exporter for metrics extracted from the Slurm resource scheduling system.

%prep
%setup -q

%build
# Empty section.

%install
rm -rf %{buildroot}
mkdir -vp %{buildroot}
mkdir -vp %{buildroot}%{_unitdir}/
mkdir -vp %{buildroot}/usr/bin
mkdir -vp %{buildroot}/usr/share/doc/prometheus-slurm-exporter-%{version}
mkdir -vp %{buildroot}/var/lib/prometheus
install -m 755 prometheus-slurm-exporter %{buildroot}/usr/bin/prometheus-slurm-exporter
install -m 644 %{SOURCE1} %{buildroot}/usr/lib/systemd/system/prometheus-slurm-exporter.service
install -m 644 %{SOURCE2} %{buildroot}/usr/share/doc/prometheus-slurm-exporter-%{version}/LICENSE
install -m 644 %{SOURCE3} %{buildroot}/usr/share/doc/prometheus-slurm-exporter-%{version}/README.md

%clean
rm -rf %{buildroot}
 
%pre
getent group prometheus >/dev/null || groupadd -r prometheus
getent passwd prometheus >/dev/null || \
    useradd -r -g prometheus -d /var/lib/slurm_exporter -s /sbin/nologin \
    -c "Prometheus exporter user" prometheus
exit 0

%post
systemctl enable %{name}.service
systemctl start %{name}.service

%preun
%systemd_preun %{name}.service

%postun
%systemd_postun_with_restart %{name}.service

%files
%defattr(-,root,root,-)
%doc LICENSE
%doc README.md
%{_bindir}/prometheus-slurm-exporter
%{_unitdir}/%{name}.service
%attr(755, prometheus, prometheus)/var/lib/prometheus

%changelog
* Fri Feb 16 2018 Matteo <m.dessalvi@gsi.de> - 0.3
- Fix issue #4
* Wed Jan 31 2018 Matteo <m.dessalvi@gsi.de> - 0.2
- Fix issue #3
* Mon Dec 04 2017 Matteo <m.dessalvi@gsi.de> - 0.1
- Full commit history: https://github.com/vpenso/prometheus-slurm-exporter/commits/master
