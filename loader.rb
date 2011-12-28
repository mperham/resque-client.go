require 'resque'

class Job
  @queue = :test_queue

  def self.perform(ctx)

  end
end

Resque.enqueue(Job, { 'mike' => 'rulez', 3 => 12 })
Resque.enqueue(Job, { 'mike' => %w(jim bob), 'hash' => { 'some' => 'subhash' } })
