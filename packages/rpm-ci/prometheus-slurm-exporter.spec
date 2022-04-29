%define        __spec_install_post %{nil}
%define          debug_package %{nil}
%define        __os_install_post %{_dbpath}/brp-compress

Name:           prometheus-slurm-exporter
Version:        %{_version}
Release:        1%{?dist}
Summary:        Prometheus exporter for SLURM metrics
Group:          Monitoring

License:        GPL 3.0
URL:            https://github.com/vpenso/prometheus-slurm-exporter

Source0:        https://github.com/vpenso/prometheus-slurm-exporter/archive/refs/tags/%{version}.tar.gz
Source1:        %{name}.service
Source2:        LICENSE
Source3:        README.md
Source4:        %{name}

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
go build -v -o %{_sourcedir}/%{name}

%install
rm -rf %{buildroot}
mkdir -vp %{buildroot}
mkdir -vp %{buildroot}%{_unitdir}/
mkdir -vp %{buildroot}/usr/bin
mkdir -vp %{buildroot}/usr/share/doc/%{name}-%{version}
mkdir -vp %{buildroot}/var/lib/prometheus
install -m 644 %{SOURCE1} %{buildroot}/usr/lib/systemd/system/%{name}.service
install -m 644 %{SOURCE2} %{buildroot}/usr/share/doc/%{name}-%{version}/LICENSE
install -m 644 %{SOURCE3} %{buildroot}/usr/share/doc/%{name}-%{version}/README.md
install -m 755 %{SOURCE4} %{buildroot}/usr/bin/%{name}

%clean
rm -rf %{buildroot}
 
%pre
getent group prometheus >/dev/null || groupadd -r prometheus
getent passwd prometheus >/dev/null || \
    useradd -r -g prometheus -d /var/lib/slurm_exporter -s /sbin/nologin \
    -c "Prometheus exporter user" prometheus
exit 0

%post
%systemd_post %{name}.service

%preun
%systemd_preun %{name}.service

%postun
%systemd_postun_with_restart %{name}.service

%files
%defattr(-,root,root,-)
%doc LICENSE
%doc README.md
%{_bindir}/%{name}
/usr/share/doc/%{name}-%{version}/
%{_unitdir}/%{name}.service
%attr(755, prometheus, prometheus)/var/lib/prometheus

