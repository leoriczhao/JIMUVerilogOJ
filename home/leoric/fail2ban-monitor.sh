#!/bin/bash

# Fail2ban监控脚本
echo "=== Fail2ban状态监控 ==="
echo "时间: $(date)"
echo ""

echo "=== 服务状态 ==="
sudo systemctl status fail2ban --no-pager -l
echo ""

echo "=== 活跃的jail ==="
sudo fail2ban-client status
echo ""

echo "=== SSH jail详细状态 ==="
sudo fail2ban-client status sshd
echo ""

echo "=== 当前被封禁的IP ==="
sudo iptables -L f2b-SSH -n 2>/dev/null | grep REJECT || echo "没有被封禁的IP"
echo ""

echo "=== 最近的SSH登录失败记录 ==="
sudo tail -10 /var/log/auth.log | grep "Failed password"
echo ""

echo "=== Fail2ban日志 ==="
sudo tail -5 /var/log/fail2ban.log