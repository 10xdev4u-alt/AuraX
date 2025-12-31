import { useState, useEffect } from 'react'
import { Activity, Cpu, Package, TrendingUp } from 'lucide-react'
import { deviceService, firmwareService, releaseService } from '../services/api'

export default function Dashboard() {
  const [stats, setStats] = useState({
    totalDevices: 0,
    provisionedDevices: 0,
    totalFirmware: 0,
    activeReleases: 0,
  })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchStats()
  }, [])

  const fetchStats = async () => {
    try {
      const [devicesRes, firmwareRes, releasesRes] = await Promise.all([
        deviceService.getAll(),
        firmwareService.getAll(),
        releaseService.getAll(),
      ])

      const devices = devicesRes.data.devices || []
      const provisioned = devices.filter(d => d.provisioned_at).length

      setStats({
        totalDevices: devices.length,
        provisionedDevices: provisioned,
        totalFirmware: firmwareRes.data.firmwares?.length || 0,
        activeReleases: releasesRes.data.releases?.filter(r => r.status === 'in_progress').length || 0,
      })
    } catch (error) {
      console.error('Failed to fetch stats:', error)
    } finally {
      setLoading(false)
    }
  }

  const cards = [
    {
      title: 'Total Devices',
      value: stats.totalDevices,
      icon: Cpu,
      color: 'bg-blue-500',
    },
    {
      title: 'Provisioned',
      value: stats.provisionedDevices,
      icon: Activity,
      color: 'bg-green-500',
    },
    {
      title: 'Firmware Versions',
      value: stats.totalFirmware,
      icon: Package,
      color: 'bg-purple-500',
    },
    {
      title: 'Active Releases',
      value: stats.activeReleases,
      icon: TrendingUp,
      color: 'bg-orange-500',
    },
  ]

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  return (
    <div className="px-4 sm:px-6 lg:px-8">
      <div className="sm:flex sm:items-center">
        <div className="sm:flex-auto">
          <h1 className="text-2xl font-semibold text-gray-900">Dashboard</h1>
          <p className="mt-2 text-sm text-gray-700">
            Overview of your IoT fleet management platform
          </p>
        </div>
      </div>

      <div className="mt-8 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        {cards.map((card) => (
          <div
            key={card.title}
            className="bg-white overflow-hidden shadow rounded-lg"
          >
            <div className="p-5">
              <div className="flex items-center">
                <div className={`flex-shrink-0 ${card.color} rounded-md p-3`}>
                  <card.icon className="h-6 w-6 text-white" />
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">
                      {card.title}
                    </dt>
                    <dd className="text-3xl font-semibold text-gray-900">
                      {card.value}
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="mt-8 bg-white shadow rounded-lg p-6">
        <h2 className="text-lg font-medium text-gray-900 mb-4">Fleet Health</h2>
        <div className="space-y-4">
          <div>
            <div className="flex justify-between text-sm mb-1">
              <span className="text-gray-600">Provisioned Devices</span>
              <span className="text-gray-900 font-medium">
                {stats.provisionedDevices} / {stats.totalDevices}
              </span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div
                className="bg-green-500 h-2 rounded-full"
                style={{
                  width: `${stats.totalDevices > 0 ? (stats.provisionedDevices / stats.totalDevices) * 100 : 0}%`,
                }}
              ></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
