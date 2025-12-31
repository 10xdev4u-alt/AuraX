import { useState, useEffect } from 'react'
import { Rocket, Plus, RefreshCw, AlertCircle, CheckCircle } from 'lucide-react'
import { releaseService, firmwareService } from '../services/api'

export default function Releases() {
  const [releases, setReleases] = useState([])
  const [firmwares, setFirmwares] = useState([])
  const [loading, setLoading] = useState(true)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [formData, setFormData] = useState({
    firmware_id: '',
    target_fleet: 'production',
    health_policy: 'auto-rollback',
  })

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    try {
      const [releasesRes, firmwaresRes] = await Promise.all([
        releaseService.getAll(),
        firmwareService.getAll(),
      ])
      setReleases(releasesRes.data.releases || [])
      setFirmwares(firmwaresRes.data.firmwares || [])
    } catch (error) {
      console.error('Failed to fetch data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleCreateRelease = async (e) => {
    e.preventDefault()
    try {
      await releaseService.create(formData)
      setShowCreateModal(false)
      setFormData({ firmware_id: '', target_fleet: 'production', health_policy: 'auto-rollback' })
      fetchData()
    } catch (error) {
      console.error('Failed to create release:', error)
      alert('Failed to create release')
    }
  }

  const getStatusBadge = (status) => {
    const statusConfig = {
      pending: { color: 'bg-gray-100 text-gray-800', icon: AlertCircle },
      in_progress: { color: 'bg-blue-100 text-blue-800', icon: RefreshCw },
      completed: { color: 'bg-green-100 text-green-800', icon: CheckCircle },
      rolled_back: { color: 'bg-red-100 text-red-800', icon: AlertCircle },
    }

    const config = statusConfig[status] || statusConfig.pending
    const Icon = config.icon

    return (
      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${config.color}`}>
        <Icon className="w-4 h-4 mr-1" />
        {status.replace('_', ' ').toUpperCase()}
      </span>
    )
  }

  const getStageBadge = (stage) => {
    const stageColors = {
      canary: 'bg-yellow-100 text-yellow-800',
      staging: 'bg-purple-100 text-purple-800',
      production: 'bg-blue-100 text-blue-800',
      completed: 'bg-green-100 text-green-800',
      rollback: 'bg-red-100 text-red-800',
    }

    return (
      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${stageColors[stage] || 'bg-gray-100 text-gray-800'}`}>
        {stage.toUpperCase()}
      </span>
    )
  }

  const getFirmwareVersion = (firmwareId) => {
    const firmware = firmwares.find(f => f.id === firmwareId)
    return firmware ? firmware.version : 'Unknown'
  }

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
          <h1 className="text-2xl font-semibold text-gray-900">Releases</h1>
          <p className="mt-2 text-sm text-gray-700">
            Manage OTA update releases and deployments
          </p>
        </div>
        <div className="mt-4 sm:mt-0 sm:ml-16 sm:flex-none">
          <button
            onClick={() => setShowCreateModal(true)}
            className="inline-flex items-center justify-center rounded-md border border-transparent bg-primary-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-primary-700"
          >
            <Plus className="h-4 w-4 mr-2" />
            Create Release
          </button>
        </div>
      </div>

      <div className="mt-8 flex flex-col">
        <div className="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div className="inline-block min-w-full py-2 align-middle md:px-6 lg:px-8">
            <div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 md:rounded-lg">
              <table className="min-w-full divide-y divide-gray-300">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                      Release ID
                    </th>
                    <th className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                      Firmware Version
                    </th>
                    <th className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                      Status
                    </th>
                    <th className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                      Stage
                    </th>
                    <th className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                      Target Fleet
                    </th>
                    <th className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                      Created At
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200 bg-white">
                  {releases.map((release) => (
                    <tr key={release.id} className="hover:bg-gray-50">
                      <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-900 font-mono">
                        {release.id.substring(0, 8)}...
                      </td>
                      <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-900 font-semibold">
                        v{getFirmwareVersion(release.firmware_id)}
                      </td>
                      <td className="whitespace-nowrap px-3 py-4 text-sm">
                        {getStatusBadge(release.status)}
                      </td>
                      <td className="whitespace-nowrap px-3 py-4 text-sm">
                        {getStageBadge(release.stage)}
                      </td>
                      <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                        {release.target_fleet || 'All'}
                      </td>
                      <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                        {new Date(release.created_at).toLocaleString()}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      {showCreateModal && (
        <div className="fixed z-10 inset-0 overflow-y-auto">
          <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={() => setShowCreateModal(false)}></div>
            <div className="inline-block align-bottom bg-white rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full sm:p-6">
              <form onSubmit={handleCreateRelease}>
                <div>
                  <div className="flex items-center">
                    <div className="flex-shrink-0 bg-primary-100 rounded-md p-3">
                      <Rocket className="h-6 w-6 text-primary-600" />
                    </div>
                    <h3 className="ml-4 text-lg leading-6 font-medium text-gray-900">
                      Create New Release
                    </h3>
                  </div>
                  <div className="mt-4 space-y-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700">
                        Firmware Version
                      </label>
                      <select
                        required
                        value={formData.firmware_id}
                        onChange={(e) => setFormData({ ...formData, firmware_id: e.target.value })}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-primary-500 focus:border-primary-500"
                      >
                        <option value="">Select firmware...</option>
                        {firmwares.map((firmware) => (
                          <option key={firmware.id} value={firmware.id}>
                            v{firmware.version} - {firmware.description}
                          </option>
                        ))}
                      </select>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700">
                        Target Fleet
                      </label>
                      <input
                        type="text"
                        value={formData.target_fleet}
                        onChange={(e) => setFormData({ ...formData, target_fleet: e.target.value })}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-primary-500 focus:border-primary-500"
                        placeholder="production"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700">
                        Health Policy
                      </label>
                      <select
                        value={formData.health_policy}
                        onChange={(e) => setFormData({ ...formData, health_policy: e.target.value })}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-primary-500 focus:border-primary-500"
                      >
                        <option value="auto-rollback">Auto Rollback</option>
                        <option value="manual">Manual</option>
                      </select>
                    </div>
                  </div>
                </div>
                <div className="mt-5 sm:mt-6 sm:grid sm:grid-cols-2 sm:gap-3 sm:grid-flow-row-dense">
                  <button
                    type="submit"
                    className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-primary-600 text-base font-medium text-white hover:bg-primary-700 focus:outline-none sm:col-start-2 sm:text-sm"
                  >
                    Create Release
                  </button>
                  <button
                    type="button"
                    onClick={() => setShowCreateModal(false)}
                    className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none sm:mt-0 sm:col-start-1 sm:text-sm"
                  >
                    Cancel
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
